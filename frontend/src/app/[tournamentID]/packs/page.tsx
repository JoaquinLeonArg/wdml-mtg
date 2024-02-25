"use client"

import { CardDisplay, CardImage } from "@/components/card";
import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { Pack, PackList } from "@/components/packlist";
import { ApiGetRequest } from "@/requests/requests";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function PacksPage(props: any) {
  let [packs, setPacks] = useState<Pack[]>([])
  let [packsError, setPacksError] = useState<string>("")

  useEffect(() => {
    ApiGetRequest({
      route: "/tournament_player",
      errorHandler: (err) => { setPacksError(err) },
      responseHandler: (res) => {
        console.log(res)
      }
    })
  }, [])

  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        <Header title="Open packs" />
        <div className="flex flex-row gap-2 justify-center">
          <PackList packs={packs} />
          <CardDisplay cardImageURLs={[
            {
              cardImageURL: "https://cards.scryfall.io/png/front/3/4/343d01cf-9806-4c2d-a993-ddc9ed248d7f.png",
              cardRarity: "rare"
            },
            {
              cardImageURL: "https://cards.scryfall.io/large/front/4/3/434515bf-de57-4c00-b0b4-c9579cc1b84c.jpg",
              cardRarity: "uncommon"
            },
            {
              cardImageURL: "https://cards.scryfall.io/large/front/4/3/434515bf-de57-4c00-b0b4-c9579cc1b84c.jpg",
              cardRarity: "uncommon"
            },
            {
              cardImageURL: "https://cards.scryfall.io/png/front/9/a/9afe8b9e-bb14-44d5-b5da-627835ee457f.png",
              cardRarity: "common"
            },
            {
              cardImageURL: "https://cards.scryfall.io/png/front/9/a/9afe8b9e-bb14-44d5-b5da-627835ee457f.png",
              cardRarity: "common"
            },
            {
              cardImageURL: "https://cards.scryfall.io/png/front/9/a/9afe8b9e-bb14-44d5-b5da-627835ee457f.png",
              cardRarity: "common"
            },
            {
              cardImageURL: "https://cards.scryfall.io/large/front/c/3/c3f1f41e-98fc-4f6b-b287-c8899dff8ab0.jpg?1562563557",
              cardRarity: "common"
            },

          ]} />
        </div>
      </div >
    </Layout >
  )
}