"use client"

import { CardDisplaySpoiler, CardFullProps } from "@/components/collectioncard"
import { rarities, mtgColors, CollectionFilter } from "@/components/collectionfilter"
import { DecklistCardListProps, DecklistCardProps, DecklistList, DecklistListProps } from "@/components/decklistcard"
import { Header, MiniHeader } from "@/components/header"
import Layout from "@/components/layout"
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests"
import { OwnedCard } from "@/types/card"
import { Deck, DeckCard } from "@/types/deck"
import { Button, ButtonGroup, Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, Pagination, Spinner } from "@nextui-org/react"
import { useEffect, useState } from "react"
import Image from "next/image"
import { groupCardsByType, groupCardsByColor, groupCardsByMV } from "./groups"
import { FormatDeck } from "./export"


export default function EditDeckPage(props: any) {
  let [deck, setDeck] = useState<Deck>()
  let [deckLoading, setDeckLoading] = useState<boolean>(false)
  let [addCardsModalOpen, setAddCardsModalOpen] = useState<boolean>(false)
  let [cardCounts, setCardCounts] = useState<{ [owned_card_id: string]: number }>({})
  let [cards, setCards] = useState<OwnedCard[]>([])
  let [preview, setPreview] = useState<string>("/cardback.webp")
  let [groupBy, setGroupBy] = useState<string>("Type")
  let [orderBy, setOrderBy] = useState<string>("Name")

  let refreshData = () => {
    ApiGetRequest({
      route: "/deck",
      query: { deck_id: props.params.deckID },
      errorHandler: (err) => {
        setDeckLoading(false)
      },
      responseHandler: (res: { deck: Deck, cards: OwnedCard[] }) => {
        setDeckLoading(false)
        setDeck(res.deck)
        let newCardCounts: { [owned_card_id: string]: number } = {}
        res.deck.cards.forEach((card) => {
          newCardCounts[card.owned_card_id] = card.count
        })
        setCardCounts(newCardCounts)
        setCards(res.cards)
      }
    })
  }

  useEffect(() => {
    setDeckLoading(true)
    refreshData()
  }, [props.params.tournamentID])

  return (
    <>
      <AddCardsModal tournamentID={props.params.tournamentID} deckID={props.params.deckID} deck={deck} cardCounts={cardCounts} isOpen={addCardsModalOpen} closeFn={() => setAddCardsModalOpen(false)} refreshDeckFn={refreshData} />
      <Layout tournamentID={props.params.tournamentID} >
        <div className="mx-16 my-16">
          <Header title={deck?.name || ""}
            endContent={deck &&
              <div className="flex flex-row gap-2">
                <Button onClick={() => deck ? navigator.clipboard.writeText(FormatDeck(deck, cards)) : null} color="warning">Copy deck to clipboard</Button>
                <Button isIconOnly onClick={() => setAddCardsModalOpen(true)} color="success">+</Button>
              </div>
            } />
          <div className="flex flex-row gap-16">
            <Image
              className="duration-75 border-2 border-white rounded-xl w-72 h-96"
              src={preview}
              alt="back"
              width={1024}
              height={768}
              quality={100}
            />
            <div className="flex flex-col w-full">
              <div className="flex flex-row items-center gap-2 mb-4">
                <p className="text-gray-500">Group by</p>
                <Dropdown>
                  <DropdownTrigger>
                    <Button
                      variant="bordered"
                      className="capitalize min-w-32"
                    >
                      {groupBy}
                    </Button>
                  </DropdownTrigger>
                  <DropdownMenu
                    className="text-white"
                  >
                    {
                      Object.keys(categoryGroupers).map((key) =>
                        <DropdownItem onPress={() => setGroupBy(key)} key={key}>{key}</DropdownItem>
                      )
                    }
                  </DropdownMenu>
                </Dropdown>
                <p className="text-gray-500 pl-10">Order by</p>
                <Dropdown>
                  <DropdownTrigger>
                    <Button
                      variant="bordered"
                      className="capitalize min-w-32"
                    >
                      {orderBy}
                    </Button>
                  </DropdownTrigger>
                  <DropdownMenu
                    className="text-white"
                  >
                    {
                      Object.keys(orders).map((key) =>
                        <DropdownItem onPress={() => setOrderBy(key)} key={key}>{key}</DropdownItem>
                      )
                    }
                  </DropdownMenu>
                </Dropdown>
              </div>
              {
                deckLoading ? <div className="flex justify-center"> <Spinner /></div> :
                  <div className="flex flex-col w-full">
                    <MiniHeader title="Main Deck"
                      endContent={
                        <div className="text-gray-300 self-end">
                          {deck?.cards.filter(card => card.board == "b_mainboard").map((card) => card.count).reduce((pv, cv) => pv + cv, 0)}
                        </div>
                      }
                    />
                    {deck ? <DeckDisplayFull
                      refreshFn={refreshData}
                      deckId={props.params.deckID}
                      cards={cards}
                      groupBy={groupBy}
                      orderBy={orderBy}
                      board="b_mainboard"
                      setPreview={setPreview}
                      deckCards={deck.cards.filter(card => card.board == "b_mainboard")}
                    /> : ""}
                    <MiniHeader title="Sideboard"
                      endContent={
                        <div className="text-gray-300 self-end">
                          {deck?.cards.filter(card => card.board == "b_sideboard").map((card) => card.count).reduce((pv, cv) => pv + cv, 0)}
                        </div>
                      }
                    />
                    {deck ? <DeckDisplayFull
                      refreshFn={refreshData}
                      deckId={props.params.deckID}
                      cards={cards}
                      groupBy={groupBy}
                      orderBy={orderBy}
                      board="b_sideboard"
                      setPreview={setPreview}
                      deckCards={deck.cards.filter(card => card.board == "b_sideboard")}
                    /> : ""}
                    <MiniHeader title="Considering"
                      endContent={
                        <div className="text-gray-300 self-end">
                          {deck?.cards.filter(card => card.board == "b_maybeboard").map((card) => card.count).reduce((pv, cv) => pv + cv, 0)}
                        </div>
                      }
                    />
                    {deck ? <DeckDisplayFull
                      refreshFn={refreshData}
                      deckId={props.params.deckID}
                      cards={cards}
                      groupBy={groupBy}
                      orderBy={orderBy}
                      board="b_maybeboard"
                      setPreview={setPreview}
                      deckCards={deck.cards.filter(card => card.board == "b_maybeboard")}
                    /> : ""}
                  </div>
              }
            </div>
          </div>
        </div>
      </Layout >
    </>
  )
}

type AddCardsModalProps = {
  tournamentID: string
  deckID: string
  deck?: Deck
  cardCounts: { [owned_card_id: string]: number }
  isOpen: boolean
  closeFn: () => void
  refreshDeckFn: () => void
}

function AddCardsModal(props: AddCardsModalProps) {
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

  let [currentCards, setCurrentCards] = useState<CardFullProps[]>([])
  let [totalPages, setTotalPages] = useState<number>(0)
  let [page, setPage] = useState<number>(1)

  let [board, setBoard] = useState<string>("b_mainboard")

  let [error, setError] = useState<string>("")

  let refreshCollection = () => {
    ApiGetRequest({
      route: "/collection",
      query: {
        filters: `name=${cardName}+tags=${tags}+rarity=${rarity}+color=${colors.W ? "W" : ""}${colors.U ? "U" : ""}${colors.B ? "B" : ""}${colors.R ? "R" : ""}${colors.G ? "G" : ""}${colors.C ? "C" : ""}+types=${types}+oracle=${oracle}+setcode=${setCode}+mv=${mv}`,
        tournament_id: props.tournamentID,
        count: 75,
        page,
      },
      errorHandler: (err) => { setError(err) },
      responseHandler: (res: { cards: OwnedCard[], count: number, max_page: number }) => {
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
            onClickFn: () => addCard(card.id),
            count: card.count,
            disabled: card.count <= props.cardCounts[card.id] || (props.cardCounts[card.id] >= 4 && !card.card_data.types.includes("Basic")),
            inDeck: props.cardCounts[card.id]
          }
        })
        setCurrentCards(cardsFull)
      }
    })
  }

  let addCard = (card_id: string) => {
    ApiPostRequest({
      route: "/deck/card",
      body: {
        owned_card_id: card_id,
        deck_id: props.deckID,
        amount: 1,
        board: board
      },
      errorHandler: (err) => { setError(err) },
      responseHandler: () => {
        props.refreshDeckFn()
        refreshCollection()
      }
    })
  }

  useEffect(() => {
    refreshCollection()
  }, [cardName, tags, rarity, colors, types, oracle, setCode, mv, page, props])

  if (!props.deck) {
    return
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
        <ModalHeader className="flex flex-col text-white gap-2">
          Add cards
          <CollectionFilter
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

        </ModalHeader>
        <ModalBody className="py-8">
          <CardDisplaySpoiler cards={currentCards} />
        </ModalBody>
        <ModalFooter>
          <div className="w-full flex flex-row gap-2 items-center justify-between">
            <div className="flex flex-row items-center gap-2">
              <div className="text-white">Add cards to</div>
              <ButtonGroup>
                <Button
                  className={`${board == "b_mainboard" ? "bg-gray-400" : ""}`}
                  onClick={() => setBoard("b_mainboard")}>
                  Main Deck
                </Button>
                <Button
                  className={`${board == "b_sideboard" ? "bg-gray-400" : ""}`}
                  onClick={() => setBoard("b_sideboard")}>
                  Sideboard
                </Button>
                <Button
                  className={`${board == "b_maybeboard" ? "bg-gray-400" : ""}`}
                  onClick={() => setBoard("b_maybeboard")}>
                  Considering
                </Button>
              </ButtonGroup>
              <Pagination onChange={(page) => setPage(page)} isCompact showControls total={totalPages} initialPage={1}></Pagination>
            </div>
            <Button color="danger" variant="flat" onPress={props.closeFn}>
              Close
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

type DeckDisplayListProps = {
  cards: OwnedCard[]
  deckCards: DeckCard[]
  setPreview: (image_url: string) => void
  groupBy: string
  orderBy: string
  board: string
  deckId: string
  refreshFn: () => void
}

let categoryGroupers: { [groupingType: string]: (deck_cards: DeckCard[], cards: { [card_id: string]: OwnedCard }) => { [category: string]: DecklistCardProps[] } } = {
  "Type": groupCardsByType,
  "Color": groupCardsByColor,
  "Mana Value": groupCardsByMV
}

let orders: { [orderType: string]: (a: DecklistCardProps, b: DecklistCardProps) => number } = {
  "Name": (a, b) => a.card.card_data.name < b.card.card_data.name ? -1 : 1,
  "Mana Value": (a, b) => a.card.card_data.mana_value < b.card.card_data.mana_value ? -1 : 1
}

function DeckDisplayFull(props: DeckDisplayListProps) {
  let [cardsByCategory, setCardsByCategory] = useState<{ [category: string]: DecklistCardProps[] }>({})

  useEffect(() => {
    if (!props.cards) return
    let cardsById: { [id: string]: OwnedCard } = {}
    props.cards.forEach((card) => {
      cardsById[card.id] = card
    })
    let grouped = categoryGroupers[props.groupBy] ? categoryGroupers[props.groupBy](props.deckCards, cardsById) : {}
    // Set common data
    Object.keys(grouped).forEach((category) => {
      grouped[category].forEach((card) => {
        card.onHoverFn = () => props.setPreview(cardsById[card.card.id].card_data.image_url)
        card.board = props.board
        card.deckId = props.deckId
        card.refreshDeckFn = props.refreshFn
      })
      grouped[category].sort((a, b) => orders[props.orderBy](a, b))
    })
    setCardsByCategory(grouped)
  }, [props])

  return (
    <div className="flex flex-row">
      <div className={`flex flex-col flex-wrap ${props.deckCards.length < 40 ? "max-h-[600px]" : "max-h-[900px]"}`}>
        {Object.keys(cardsByCategory).map((category) => (<DecklistList key={category} cards={cardsByCategory[category]} category={category} />))}
      </div>
    </div>
  )
}

