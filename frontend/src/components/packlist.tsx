import { ApiGetRequest, ApiPostRequest } from "@/requests/requests"
import { CardData } from "@/types/card"
import { BoosterPack, BoosterPackData } from "@/types/tournamentPlayer"
import { useState } from "react"

export type PackListProps = {
  packs: BoosterPack[]
  openPackHandler: (cards: CardData[]) => void
  tournamentID: string
}

export function PackList(props: PackListProps) {
  return (
    <div className="flex flex-col gap-2">
      {
        props.packs.map((pack: BoosterPack, i: number) => {
          return (
            <BoosterPack key={i} pack={pack} openPackHandler={props.openPackHandler} tournamentID={props.tournamentID} />
          )
        })
      }
    </div>
  )
}

type BoosterPackProps = {
  pack: BoosterPack
  openPackHandler: (cards: CardData[]) => void
  tournamentID: string
}

function BoosterPack(props: BoosterPackProps) {
  let [error, setError] = useState<string>("")
  let sendOpenPackRequest = () => ApiPostRequest({
    route: `/boosterpacks/open`,
    query: { tournamentID: props.tournamentID },
    body: {
      booster_pack_data: props.pack.data
    },
    errorHandler: (err) => { setError(err) },
    responseHandler: (res: { card_data: CardData[] }) => {
      console.log(res)
      props.openPackHandler(res.card_data)
    }
  })

  return (
    ""
  )
}