"use client"

import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests";
import { OwnedCard } from "@/types/card";
import { useCallback, useEffect, useState } from "react";
import { Input, Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Button, ButtonGroup, Pagination, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, Textarea } from "@nextui-org/react";
import Image from "next/image"
import { CardDisplaySpoiler, CardFullProps } from "@/components/collectioncard";
import { CollectionFilter, mtgColors, rarities } from "@/components/collectionfilter";
import { useDropzone } from "react-dropzone";
import { parse } from "csv-parse/sync";



export default function CollectionPage(props: any) {
  // Filters
  let [cardName, setCardName] = useState<string>("")
  let [tags, setTags] = useState<string>("")
  let [rarity, setRarity] = useState<keyof typeof rarities>("")
  let [colors, setColors] = useState<mtgColors>
    ({ W: false, U: false, B: false, R: false, G: false, C: false })
  let [types, setTypes] = useState<string>("")
  let [oracle, setOracle] = useState<string>("")
  let [setCode, setSetCode] = useState<string>("")
  let [mv, setMv] = useState<string>("")

  let [totalResults, setTotalResults] = useState<number>(0)
  let [totalPages, setTotalPages] = useState<number>(0)
  let [page, setPage] = useState<number>(1)

  let [error, setError] = useState<string>("")
  let [currentCards, setCurrentCards] = useState<CardFullProps[]>([])
  let [isImportModalOpen, setIsImportModalOpen] = useState<boolean>(false)

  useEffect(() => {
    refreshData()
  }, [cardName, tags, rarity, colors, types, oracle, setCode, mv, page, props.params.tournamentID])

  let refreshData = () => {
    ApiGetRequest({
      route: "/collection",
      query: {
        filters: `name=${cardName}+tags=${tags}+rarity=${rarity}+color=${colors.W ? "W" : ""}${colors.U ? "U" : ""}${colors.B ? "B" : ""}${colors.R ? "R" : ""}${colors.G ? "G" : ""}${colors.C ? "C" : ""}+types=${types}+oracle=${oracle}+setcode=${setCode}+mv=${mv}`,
        tournament_id: props.params.tournamentID,
        count: 75,
        page,
      },
      errorHandler: (err) => { setError(err) },
      responseHandler: (res: { cards: OwnedCard[], count: number, max_page: number }) => {
        setTotalResults(res.count)
        setTotalPages(res.max_page)
        if (!res.cards || res.cards.length == 0) {
          setCurrentCards([])
        }
        // Workaround to have the flip callback see the card list
        let cardsFull = res.cards.map((card: OwnedCard) => {
          return {
            card: card.card_data,
            flipped: true,
            showRarityWhenFlipped: true,
            onClickFn: () => { },
            count: card.count
          }
        })
        cardsFull.forEach((card: CardFullProps, index: number) => { card.onClickFn = () => flipCard(index, cardsFull) })
        setCurrentCards(cardsFull)
      }
    })
  }

  let flipCard = (index: number, oldCards: CardFullProps[]) => {
    let cards = [...oldCards]
    if (cards[index].card.back_image_url != "") {
      cards[index].flipped = !cards[index].flipped
    } else {
      cards[index].flipped = true
    }
    setCurrentCards(cards)
  }

  return (
    <>
      <ImportCardsModal tournamentID={props.params.tournamentID} isOpen={isImportModalOpen} closeFn={() => setIsImportModalOpen(false)} refreshFn={refreshData} />
      <Layout tournamentID={props.params.tournamentID}>
        <div className="mx-16 my-16">
          <Header title="Collection" endContent={<Button onClick={() => setIsImportModalOpen(true)} color="warning">Import collection</Button>} />
          <div className="pb-4">
            <CollectionFilter
              count={totalResults}
              setCardName={setCardName}
              setTags={setTags}
              setRarity={setRarity}
              rarity={rarity}
              setColors={setColors}
              colors={colors}
              setTypes={setTypes}
              setSetCode={setSetCode}
              setOracle={setOracle}
              setMv={setMv}
            />
          </div>
          <div className="pb-4">
            <CardDisplaySpoiler cards={currentCards} />
          </div>
          <div className="flex flex-col items-center gap-4">
            <Pagination onChange={(page) => setPage(page)} isCompact showControls total={totalPages} initialPage={1}></Pagination>
          </div>
        </div>
      </Layout >
    </>
  )
}


type ImportCardsModalProps = {
  tournamentID: string
  isOpen: boolean
  closeFn: () => void
  refreshFn: () => void
}

function ImportCardsModal(props: ImportCardsModalProps) {
  let [deckName, setDeckName] = useState<string>("")
  let [deckDescription, setDeckDescription] = useState<string>("")
  let [error, setError] = useState<string>("")
  let [isLoading, setIsLoading] = useState<boolean>(false)
  let [data, setData] = useState<string>("")
  let [count, setCount] = useState<number>(0)
  let [uniqueCount, setUniqueCount] = useState<number>(0)

  let sendCreateDeckRequest = () => {
    setError("")
    setIsLoading(true)
    ApiPostRequest({
      route: "/collection/import",
      query: {
        tournament_id: props.tournamentID
      },
      body: data,
      rawBody: true,
      errorHandler: (err) => {
        setIsLoading(false)
        setError(err)
        props.refreshFn()
      },
      responseHandler: () => {
        setIsLoading(false)
        props.closeFn()
        props.refreshFn()
      }
    })
  }

  const onDrop = useCallback((acceptedFiles: Blob[]) => {
    acceptedFiles.forEach((file: Blob) => {
      const reader = new FileReader()

      reader.onabort = () => setError("File read aborted")
      reader.onerror = () => setError("File reading failed")
      reader.onload = (e) => {
        if (!e.target) {
          setError("File reading failed")
          return
        }
        setData(e.target.result as string)
        console.log(e.target.result as string)
        let currentCount = 0
        let currentUniqueCount = 0
        parse(e.target.result as string).forEach((line: string[], index: number) => {
          console.log(line)
          if (index == 0) { return }
          if (line.length != 13) {
            setError("File has unrecognized rows")
            return
          }
          let c = Number(line[0])
          if (!c || c == 0) {
            setError("File has unrecognized rows")
            return
          }
          currentCount += c
          currentUniqueCount += 1
        })
        setCount(currentCount)
        setUniqueCount(currentUniqueCount)
      }
      reader.readAsText(file)
    })

  }, [])
  const { getRootProps, getInputProps } = useDropzone({ onDrop })

  return (
    <Modal
      hideCloseButton
      onClose={props.closeFn}
      isOpen={props.isOpen}
      placement="top-center"
      size="xl"
      className="z-100"
    >
      <ModalContent>
        <ModalHeader className="flex flex-col text-white gap-1">Import cards to collection</ModalHeader>
        <ModalBody>
          <div className="text-blue-200 text-xl border-2 px-4 py-12 text-center rounded-2xl border-blue-200 cursor-pointer" {...getRootProps()}>
            <input {...getInputProps()} />
            <p>Drag .csv file here or click to browse</p>
          </div>
          {count > 0 && uniqueCount > 0 && <p className="text-sm font-light text-gray-400 h-2">Recognized {count} cards total, {uniqueCount} unique</p>}
          <p className="text-sm font-light text-yellow-400 h-2">Warning! This action cannot be undone.</p>
          <p className="text-sm font-light text-red-400 h-2">{error}</p>
        </ModalBody>
        <ModalFooter>
          <Button isDisabled={isLoading} color="danger" variant="flat" onPress={props.closeFn}>
            Cancel
          </Button>
          <Button isLoading={isLoading} color="success" onPress={sendCreateDeckRequest}>
            Import
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}