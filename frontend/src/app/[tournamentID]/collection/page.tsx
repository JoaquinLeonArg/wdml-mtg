"use client"

import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest } from "@/requests/requests";
import { OwnedCard } from "@/types/card";
import { useEffect, useState } from "react";
import { Input, Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Button, ButtonGroup, Pagination } from "@nextui-org/react";
import Image from "next/image"
import { CardDisplaySpoiler, CardFullProps } from "@/components/collectioncard";
import { CollectionFilter, mtgColors, rarities } from "@/components/collectionfilter";



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

  useEffect(() => {
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
  }, [cardName, tags, rarity, colors, types, oracle, setCode, mv, page, props.params.tournamentID])

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
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        <Header title="Collection" />
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
  )
}