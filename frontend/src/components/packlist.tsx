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
    <div className="flex flex-row items-center justify-between w-80 h-16 bg-primary-400 hover:bg-primary-200 rounded-lg" onClick={sendOpenPackRequest}>
      <div className="h-full p-4 text-base flex flex-row items-center justify-between gap-4">
        <span className="text-3xl">
          {props.pack.available}
        </span>
        <div className="flex flex-col">
          <span className="text-md">
            {props.pack.data.set_name}
          </span>
          <span className="font-bold text-xl">
            {props.pack.data.set_code}
          </span>
        </div>
      </div>
      <svg className="mr-4 w-8 h-8" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fillRule="evenodd" d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z"></path></svg>
    </div>
  )
}