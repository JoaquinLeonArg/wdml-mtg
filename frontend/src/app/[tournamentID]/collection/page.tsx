"use client"

import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests";
import { CardData, OwnedCard } from "@/types/card";
import { useCallback, useEffect, useState } from "react";
import { Input, Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Button, ButtonGroup, Pagination, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, Textarea, Spinner } from "@nextui-org/react";
import Image from "next/image"
import { CardDisplaySpoiler, CardFullProps } from "@/components/collectioncard";
import { CollectionFilter, mtgColors, rarities } from "@/components/collectionfilter";
import { useDropzone } from "react-dropzone";
import { parse } from "csv-parse/sync";
import { DoTradeUpRequest } from "@/requests/collection";



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

  let [cauldronModeActive, setCauldronModeActive] = useState<boolean>(false)
  let [cauldronCards, setCauldronCards] = useState<{ count: number, card: OwnedCard }[]>([])

  useEffect(() => {
    refreshData()
  }, [cardName, tags, rarity, colors, types, oracle, setCode, mv, page, cauldronCards, props.params.tournamentID])

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
          let foundCauldronCard = cauldronCards.find((cauldronCard) => cauldronCard.card.id == card.id)
          return {
            card: card.card_data,
            flipped: true,
            showRarityWhenFlipped: true,
            inDeckText: `${foundCauldronCard ? foundCauldronCard.count + ' in the Cauldron' : ""}`,
            count: card.count
          }
        })
        cardsFull.forEach((card: CardFullProps, index: number) => {
          if (cauldronModeActive) {
            card.onClickFn = () => {
              if (cauldronCards.reduce((prev, cauldronCard) => prev + cauldronCard.count, 0) >= 30) return
              let foundCauldronCard = cauldronCards.findIndex((cauldronCard) => cauldronCard.card.id == res.cards[index].id)
              if (foundCauldronCard == -1) {
                setCauldronCards([...cauldronCards, { count: 1, card: res.cards[index] }])
              } else {
                let newCauldronCards = [...cauldronCards]
                newCauldronCards[foundCauldronCard].count += 1
                setCauldronCards(newCauldronCards)
              }
            }
            let foundCauldronCard = cauldronCards.find((cauldronCard) => cauldronCard.card.id == res.cards[index].id)
            card.disabled = foundCauldronCard != undefined && foundCauldronCard.count >= res.cards[index].count
          } else {
            card.onClickFn = () => flipCard(index, cardsFull)
          }
        })
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
          <Header title="Collection" endContent={
            <div className="flex flex-row items-end justify-center gap-2">
              <Button onClick={() => setIsImportModalOpen(true)} color="warning">Import collection</Button>
              <Button onClick={() => {
                setCauldronModeActive(!cauldronModeActive)
                setCauldronCards([])
              }} color="success" className="w-44">Cauldron mode: {cauldronModeActive ? "ON" : "OFF"}</Button>
            </div>
          } />
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
              onChange={() => { setPage(1) }}
            />
          </div>
          <div className="pb-4">
            <CardDisplaySpoiler cards={currentCards} />
          </div>
          <div className="flex flex-col items-center gap-4">
            <Pagination onChange={(page) => setPage(page)} isCompact showControls total={totalPages} initialPage={1}></Pagination>
          </div>
          {
            cauldronModeActive &&
            <CauldronFloat
              tournamentID={props.params.tournamentID}
              cauldronCards={cauldronCards}
              setCauldronCardsFn={(cards: { count: number, card: OwnedCard }[]) => setCauldronCards(cards)}
            />
          }
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

type CauldronFloatProps = {
  tournamentID: string
  cauldronCards: { count: number, card: OwnedCard }[]
  setCauldronCardsFn: (cauldronCards: { count: number, card: OwnedCard }[]) => void
}

function CauldronFloat(props: CauldronFloatProps) {
  let [cauldronModalOpen, setCauldronModalOpen] = useState<boolean>(false)

  let totalCauldronCards = props.cauldronCards.reduce((prev, curr) => prev + curr.count, 0)



  return (
    <>
      <CauldronModal
        tournamentID={props.tournamentID}
        inputCards={props.cauldronCards}
        isOpen={cauldronModalOpen}
        closeFn={() => setCauldronModalOpen(false)}
        refreshFn={() => { }}
        setCardsFn={(cauldronCards) => props.setCauldronCardsFn(cauldronCards)}
      />
      {
        cauldronModalOpen ? "" :
          <div className="fixed flex flex-row items-center justify-between bottom-4 right-16 transform bg-gray-800 text-white px-8 py-2 rounded-full border-white border-large z-[100]">
            {
              props.cauldronCards.length == 0 ?
                "The cauldron is empty, click on a card to add it" :
                <div className="flex flex-row items-center justify-center gap-8">
                  {totalCauldronCards}/30 cards on the cauldron
                  {totalCauldronCards == 30 ? <Button onClick={() => setCauldronModalOpen(true)} color="success">Reroll cards</Button> : ""}
                </div>
            }
          </div>
      }
    </>
  )
}

type CauldronModalProps = {
  tournamentID: string
  inputCards: { count: number, card: OwnedCard }[]
  isOpen: boolean
  setCardsFn: (cauldronCards: { count: number, card: OwnedCard }[]) => void,
  closeFn: () => void
  refreshFn: () => void
}

function CauldronModal(props: CauldronModalProps) {
  let [error, setError] = useState<string>()
  let [isLoading, setIsLoading] = useState<boolean>()
  let [newCards, setNewCards] = useState<CardFullProps[]>([])
  let [cardsFull, setCardsFull] = useState<CardFullProps[]>([])

  useEffect(() => {
    console.log(props.inputCards, cardsFull)
    let cf = props.inputCards.map((cauldronCard: { count: number, card: OwnedCard }, index: number) => {
      return {
        card: cauldronCard.card.card_data,
        flipped: true,
        showRarityWhenFlipped: true,
        inDeckText: `${cauldronCard.count} in the Cauldron`,
        dropdownOptions: [
          {
            key: "remove-from-cauldron",
            label: "Remove one to the Cauldron",
            enabled: () => true,
            action: () => {
              let newCauldronCards = [...props.inputCards]
              if (cauldronCard.count == 1) {
                newCauldronCards.splice(index, 1)
                props.setCardsFn(newCauldronCards)
              } else {
                newCauldronCards[index].count -= 1
                props.setCardsFn(newCauldronCards)
              }
            }
          }
        ],
        count: 0
      }
    })
    setCardsFull(cf)
  }, [props.inputCards])

  let sendTradeupRequest = () => {
    let cardsReq: { [card_id: string]: number } = {}
    props.inputCards.forEach((card) => cardsReq[card.card.id] = card.count)
    setError("")
    setIsLoading(true)
    DoTradeUpRequest(
      cardsReq,
      props.tournamentID,
      (cards) => {
        setIsLoading(false)
        props.refreshFn()
        setNewCards(cards.map((card) => {
          return {
            card: card,
            flipped: true,
            showRarityWhenFlipped: true,
            count: 0
          }
        }))
        props.setCardsFn([])
      },
      (err) => {
        setIsLoading(false)
        setError(err)
        props.refreshFn()
      },
    )
  }

  return (
    <Modal
      hideCloseButton
      onClose={props.closeFn}
      isOpen={props.isOpen}
      placement="top-center"
      size="full"
      className="w-[95vw]"
      scrollBehavior="inside"
    >
      <ModalContent>
        <ModalHeader className="flex flex-col text-white gap-1">The Cauldron of Eternity</ModalHeader>
        <ModalBody>
          {
            isLoading ? <Spinner /> :
              newCards.length == 0 ?
                <>
                  <p className="text-sm font-light text-white">The following 30 cards will be rerolled into 5 new cards:</p>
                  <p className="text-sm font-light text-yellow-400 h-2">Warning! This action cannot be undone.</p>
                  <p className="text-sm font-light text-gray-400 h-2">The sets and rarities of the new cards are decided based on the sets and rarities of the input cards.</p>
                  <p className="text-sm font-light text-red-400 h-2">{error}</p>
                  <div className="pb-4">
                    <CardDisplaySpoiler cards={cardsFull} />
                  </div>
                </> :
                <>
                  <p className="text-sm font-black text-white">This are your new cards!</p>
                  <div className="pb-4">
                    <CardDisplaySpoiler cards={newCards} />
                  </div>
                </>
          }
        </ModalBody>
        <ModalFooter>
          <Button isDisabled={isLoading} color="danger" variant="flat" onPress={props.closeFn}>
            Cancel
          </Button>
          <Button isLoading={isLoading} color="success" onPress={sendTradeupRequest}>
            Reroll cards
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}