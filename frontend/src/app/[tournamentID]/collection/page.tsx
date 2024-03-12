"use client"

import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest } from "@/requests/requests";
import { OwnedCard } from "@/types/card";
import { useEffect, useState } from "react";
import { CardImage } from "@/components/card"
import { Input, Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Button } from "@nextui-org/react";



export default function CollectionPage(props: any) {
  let [error, setError] = useState<string>("")
  let [currentCards, setCurrentCards] = useState<OwnedCard[]>([])

  useEffect(() => {
    ApiGetRequest({
      route: "/collection",
      query: {
        filters: "color<R",
        tournamentID: props.params.tournamentID,
        count: 75,
        page: 1,
      },
      errorHandler: (err) => { setError(err) },
      responseHandler: (res: { cards: OwnedCard[] }) => {
        console.log(res)
        setCurrentCards(res.cards)
      }
    })
  }, [])

  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        <Header title="Collection" />
        <div className="flex flex-row gap-2 mb-12 items-center">
          <Input
            onChange={(e) => console.log(e)}
            variant="bordered"
            type="string"
            placeholder="Set code"
            className="text-white"
          />
          <Dropdown>
            <DropdownTrigger>
              <Button
                variant="bordered"
                className="capitalize"
              >
                {/* {selectedValue} */}
                Any rarity
              </Button>
            </DropdownTrigger>
            <DropdownMenu
              onChange={(e) => console.log(e)}
              className="text-white"

            >
              <DropdownItem key="any">Any rarity</DropdownItem>
              <DropdownItem key="common">Common</DropdownItem>
              <DropdownItem key="uncommon">Uncommon</DropdownItem>
              <DropdownItem key="rare">Rare</DropdownItem>
              <DropdownItem key="mythic">Mythic Rare</DropdownItem>
            </DropdownMenu>
          </Dropdown>
          <Button>
            <i className="ms ms-u text-mana-blue"></i>
          </Button>
        </div>
        <div className="flex flex-row gap-2 justify-between">
          {currentCards && < CollectionDisplay cards={currentCards} />}
        </div>
      </div>
    </Layout >
  )
}

export type CollectionDisplayProps = {
  cards: OwnedCard[]
}

export function CollectionDisplay(props: CollectionDisplayProps) {
  return (
    <div className="flex flex-wrap flex-row gap-2 items-center justify-center">
      {props.cards.map((card: OwnedCard, i: number) => {
        return (
          <CardImage startFaceUp card_rarity={card.card_data.card_rarity} image_url={card.card_data.image_url} key={i} />
        )
      })}
    </div>
  )
}