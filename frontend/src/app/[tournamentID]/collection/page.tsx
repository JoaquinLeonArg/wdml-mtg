"use client"

import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest } from "@/requests/requests";
import { OwnedCard } from "@/types/card";
import { useEffect, useState } from "react";
import { Input, Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Button, ButtonGroup } from "@nextui-org/react";
import Image from "next/image"
import { CardDisplaySpoiler, CardFullProps } from "@/components/card";

const rarities = {
  "": "Any rarity",
  "common": "Common",
  "uncommon": "Uncommon",
  "rare": "Rare",
  "mythic": "Mythic"
}

export default function CollectionPage(props: any) {
  // Filters
  let [cardName, setCardName] = useState<string>("")
  let [tags, setTags] = useState<string>("")
  let [rarity, setRarity] = useState<keyof typeof rarities>("")
  let [colors, setColors] = useState<{ W: boolean, U: boolean, B: boolean, R: boolean, G: boolean, C: boolean }>
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
        <div className="flex flex-row gap-2 mb-2 items-center">
          <Input
            onChange={(e) => setCardName(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Card name"
            className="text-white max-w-96"
          />
          <Input
            onChange={(e) => setTags(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Tags"
            className="text-white max-w-96"
          />
          <Dropdown>
            <DropdownTrigger>
              <Button
                variant="bordered"
                className="capitalize"
              >
                {rarities[rarity]}
              </Button>
            </DropdownTrigger>
            <DropdownMenu
              className="text-white"
            >
              {
                Object.keys(rarities).map((key) =>
                  <DropdownItem onPress={() => setRarity(key as keyof typeof rarities)} key={key}>{rarities[key as keyof typeof rarities]}</DropdownItem>
                )
              }
            </DropdownMenu>
          </Dropdown>
          <ButtonGroup variant="ghost" fullWidth={false}>
            {
              Object.keys(colors).map((key) =>
                <Button isIconOnly
                  className={`text-xl ${colors[key as keyof typeof colors] && "bg-gray-500"}`}
                  variant="ghost"
                  onClick={() => setColors({ ...colors, [key as keyof typeof colors]: !colors[key as keyof typeof colors] })}
                  key={key}
                >
                  <i className={`ms ms-${key.toLowerCase()} ms-cost`}></i>
                </Button>
              )
            }
          </ButtonGroup>
        </div>
        <div className="flex flex-row gap-2 mb-2 items-center">
          <Input
            onChange={(e) => setTypes(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Types"
            className="text-white max-w-96"
          />
          <Input
            onChange={(e) => setOracle(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Oracle"
            className="text-white max-w-96"
          />
          <Input
            onChange={(e) => setSetCode(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Set code"
            className="text-white max-w-24"
          />
          <Input
            onChange={(e) => setMv(e.target.value)}
            variant="bordered"
            type="number"
            placeholder="Any"
            className="text-white max-w-28"
            endContent={
              <div className="pointer-events-none flex items-center">
                <span className="text-gray-300 text-small">MV</span>
              </div>
            }
          />
        </div>
        <div className="text-white my-6">
          Found: {totalResults} cards
        </div>
        <CardDisplaySpoiler cards={currentCards} />
        <ButtonGroup variant="ghost" className="flex flex-row items-center mt-8">
          <Button isIconOnly
            className="text-xl bg-gray-500"
            variant="ghost"
            isDisabled={page == 1}
            onClick={() => { if (page > 1) setPage(page - 1) }}
          >
            {"<"}
          </Button>
          <Button isIconOnly
            className="text-xl bg-gray-500"
            variant="ghost"
            isDisabled={page >= totalPages}
            onClick={() => { if (page < totalPages) setPage(page + 1) }}
          >
            {">"}
          </Button>
        </ButtonGroup>
      </div>
    </Layout >
  )
}