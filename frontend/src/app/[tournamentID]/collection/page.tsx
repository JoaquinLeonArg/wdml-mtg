"use client"

import { CardDisplay } from "@/components/card";
import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { PackList } from "@/components/packlist";
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests";
import { CardData, OwnedCard } from "@/types/card";
import { BoosterPack, BoosterPackData, TournamentPlayer } from "@/types/tournamentPlayer";
import { Autocomplete, AutocompleteItem, Button, Input } from "@nextui-org/react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { BsFillTrashFill } from "react-icons/bs";
import { CardImage } from "@/components/card"



export default function CollectionPage(props: any) {
  let router = useRouter()
  let [tournamentPlayer, setTournamentPlayer] = useState<TournamentPlayer>()
  let [error, setError] = useState<string>("")
  let [currentCards, setCurrentCards] = useState<OwnedCard[]>([])

  useEffect(() => {
    ApiGetRequest({
      route: "/collection",
      query: {
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
        <div className="flex flex-row gap-2 justify-between">
          <CollectionDisplay cards={currentCards} />
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