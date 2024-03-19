"use client"

import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests";
import { CardData } from "@/types/card";
import { BoosterPackData, TournamentPlayer } from "@/types/tournamentPlayer";
import { Autocomplete, AutocompleteItem, Button, ButtonGroup, Input, Listbox, ListboxItem, Spinner } from "@nextui-org/react";
import { useEffect, useState } from "react";
import { BsFillTrashFill } from "react-icons/bs";
import { CardDisplaySpoiler, CardFullProps } from "@/components/collectioncard";

export default function PacksPage(props: any) {
  let [tournamentPlayer, setTournamentPlayer] = useState<TournamentPlayer>()
  let [currentCards, setCurrentCards] = useState<CardFullProps[]>([])
  let [boostersLoading, setBoostersLoading] = useState<boolean>(false)
  let [boostersVisible, setBoostersVisible] = useState<boolean>(true)
  let [flipAllCards, setFlipAllCards] = useState<boolean>(false)
  let [flipAllCurrentIndex, setFlipAllCurrentIndex] = useState<number>(0)


  useEffect(() => {
    refreshData()
  }, [props.params.tournamentID])

  useEffect(() => {
    const interval = setInterval(() => {
      if (!flipAllCards) { setFlipAllCurrentIndex(0); return }
      if (flipAllCurrentIndex < currentCards.length) {
        flipCard(flipAllCurrentIndex, currentCards)
      } else {
        setFlipAllCards(false)
      }
      setFlipAllCurrentIndex(flipAllCurrentIndex + 1)
    }, 75)
    return () => clearInterval(interval)
  })

  let flipCard = (index: number, oldCards: CardFullProps[]) => {
    let cards = [...oldCards]
    if (cards[index].card.back_image_url != "") {
      cards[index].flipped = !cards[index].flipped
    } else {
      cards[index].flipped = true
    }
    setCurrentCards(cards)
  }


  let refreshData = () => {
    setBoostersLoading(true)
    ApiGetRequest({
      route: "/tournament_player/tournament",
      query: { tournament_id: props.params.tournamentID },
      errorHandler: (err) => {
        setBoostersLoading(false)
      },
      responseHandler: (res: { tournament_player: TournamentPlayer }) => {
        setBoostersLoading(false)
        setTournamentPlayer(res.tournament_player)
      }
    })
  }

  let sendOpenPackRequest = (booster_pack_data: BoosterPackData) => {
    setBoostersVisible(false)
    setBoostersLoading(true)
    ApiPostRequest({
      route: `/boosterpacks/open`,
      query: { tournament_id: props.params.tournamentID },
      body: {
        booster_pack_data
      },
      errorHandler: (err) => {
        setBoostersVisible(true)
        setBoostersLoading(false)
      },
      responseHandler: (res: { card_data: CardData[] }) => {
        console.log(res)
        setBoostersLoading(false)
        // Workaround to have the flip callback see the card list
        let cardsFull = res.card_data.map((cardData: CardData) => ({
          card: cardData,
          flipped: false,
          onClickFn: () => { },
          showRarityWhenFlipped: true
        }))
        cardsFull.forEach((card: CardFullProps, index: number) => { card.onClickFn = () => flipCard(index, cardsFull) })
        setCurrentCards(cardsFull)
      }
    })
  }

  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        {/* Open packs admin */}
        {/* TODO: Move this to an admin page */}
        {
          tournamentPlayer?.access_level == "al_administrator" &&
          (
            <>
              <Header title="Add packs" />
              <div className="flex flex-row gap-2 justify-left">
                <AddPacks tournamentID={props.params.tournamentID} refreshFunc={refreshData} />
              </div>
            </>
          )
        }
        {/* Select what pack to open */}
        <Header title="Open packs" />
        {boostersLoading || !tournamentPlayer ? <div className="flex justify-center"> <Spinner /></div> :
          <div className="flex flex-row gap-2 justify-center">
            {boostersVisible && (
              <div className="bg-gray-800 w-[450px] border-small px-1 py-2 rounded-small border-default-200">
                <Listbox>
                  {
                    tournamentPlayer.game_resources.booster_packs.map((booster_pack) =>
                      <ListboxItem
                        className="text-white"
                        onPress={() => sendOpenPackRequest(booster_pack.data)}
                        key={booster_pack.data.set_code}
                        startContent={<div className="text-gray-500 w-16">{booster_pack.data.set_code}</div>}
                        endContent={<div className="text-gray-500 text-right">{`(${booster_pack.available})`}</div>}
                      >
                        {`${booster_pack.data.set_code} - ${booster_pack.data.set_name}`}
                      </ListboxItem>
                    )
                  }

                </Listbox>
              </div>
            )}
            {/* Show opened booster contents */}
            {currentCards.length && (
              <div className="flex flex-col gap-8 justify-center">
                <ButtonGroup>
                  <Button color="success" onClick={() => setFlipAllCards(true)}>Flip all</Button>
                  <Button color="danger" onClick={() => { setCurrentCards([]); refreshData(); setBoostersVisible(true) }}>Close</Button>
                </ButtonGroup>
                <CardDisplaySpoiler cards={currentCards} />
              </div>
            )}
          </div>
        }
      </div >
    </Layout >
  )
}

type AddPacksProps = {
  tournamentID: string
  refreshFunc: () => void
}

type AddedPacks = {
  count: number
  set: string
  type: string
}[]

function AddPacks(props: AddPacksProps) {
  let [availablePacks, setAvailablePacks] = useState<BoosterPackData[]>([])
  let [error, setError] = useState<string>("")
  let [addedPacks, setAddedPacks] = useState<AddedPacks>([{ type: "bt_draft", count: 0, set: "" }])
  let [loading, setLoading] = useState<boolean>(false)

  let isButtonDisabled = addedPacks.length == 0 || addedPacks[0].count <= 0 || !addedPacks[0].set

  let sendAddPacksRequest = () => {
    setError("")
    setLoading(true)
    ApiPostRequest({
      body: {
        booster_packs: addedPacks
      },
      route: "/boosterpacks",
      query: { tournament_id: props.tournamentID },
      responseHandler: (_) => {
        setLoading(false)
        props.refreshFunc()
      },
      errorHandler: (err) => {
        setLoading(false)
        switch (err) {
          case "INVALID_AUTH":
            setError("Invalid credentials")
          case "NO_DATA":
            setError("Invalid number or type of packs")
          default:
            setError("An error ocurred")
        }
      }
    })
  }

  useEffect(() => {
    refreshAvailableBoosters()
  }, [props.tournamentID])

  let refreshAvailableBoosters = () => {
    ApiGetRequest({
      route: "/boosterpacks",
      query: { tournament_id: props.tournamentID },
      responseHandler: (res: { booster_packs: BoosterPackData[] }) => {
        setAvailablePacks(res.booster_packs)
      },
      errorHandler: (err) => {
        setError(err)
      }
    })
  }

  return (
    <div>
      <form onSubmit={sendAddPacksRequest} className="flex flex-col gap-2 items-center">
        <div className="flex flex-row gap-2 items-center">
          <Input
            disabled={loading}
            onChange={(e) => setAddedPacks([{ ...addedPacks[0], count: Number(e.target.value) }])}
            variant="bordered"
            type="number"
            min={0}
            label="Amount"
            placeholder="Number of packs"
            labelPlacement="inside"
            className="text-white"
            endContent={
              <div className="pointer-events-none flex items-center">
                <span className="text-gray-300 text-small">packs</span>
              </div>
            }
          />
          <Autocomplete
            disabled={loading}
            onInputChange={(value) => setAddedPacks([{ ...addedPacks[0], set: availablePacks.find(v => `${v.set_code} - ${v.set_name}` == value)?.set_code || "" }])}
            id="set"
            label="Booster pack"
            labelPlacement="inside"
            placeholder="Select a boster pack"
            className="text-white max-w-xs"
            defaultItems={availablePacks.map((val) => { return { value: val.set_code, label: `${val.set_code} - ${val.set_name}` } })}
          >
            {(item) => <AutocompleteItem className="text-white" key={item.value}>{item.label}</AutocompleteItem>}
          </Autocomplete>
          <Button size="md" isIconOnly aria-label="Like" disabled>
            <BsFillTrashFill className="w-6 h-6" />
          </Button>
        </div>
        <Button isLoading={loading} onPress={sendAddPacksRequest} disabled={isButtonDisabled} size="md" className={`w-full ${isButtonDisabled || "bg-primary-500"}`} aria-label="Like">Add packs</Button>
      </form >
      <p className="text-sm font-light text-red-400">{error}</p>
    </div >
  )
}

