"use client"

import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests";
import { CardData } from "@/types/card";
import { BoosterPack, TournamentPlayer } from "@/types/tournamentPlayer";
import { BoosterPack as GlobalBoosterPack } from "@/types/boosterPack";
import { Autocomplete, AutocompleteItem, Button, ButtonGroup, Input, Listbox, ListboxItem, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, Spinner } from "@nextui-org/react";
import { useEffect, useState } from "react";
import { BsFillCartFill, BsHypnotize } from "react-icons/bs";
import { CardDisplaySpoiler, CardFullProps } from "@/components/collectioncard";
import { User } from "@/types/user";
import { DoBuyBoosterPackRequest, DoGetAvailableBoosterPacksRequest } from "@/requests/boosterpacks";
import { DoGetTournamentStoreRequest } from "@/requests/tournament";
import { Store } from "@/types/tournament";

export default function PacksPage(props: any) {
  let [tournamentPlayer, setTournamentPlayer] = useState<TournamentPlayer>()
  let [currentCards, setCurrentCards] = useState<CardFullProps[]>([])
  let [boostersLoading, setBoostersLoading] = useState<boolean>(false)
  let [boostersVisible, setBoostersVisible] = useState<boolean>(true)
  let [flipAllCards, setFlipAllCards] = useState<boolean>(false)
  let [flipAllCurrentIndex, setFlipAllCurrentIndex] = useState<number>(0)
  let [lastOpenedPack, setLastOpenedPack] = useState<BoosterPack>()
  let [isPackStoreModalOpen, setIsPackStoreModalOpen] = useState<boolean>(false)

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

  let sendOpenPackRequest = (booster_pack: BoosterPack) => {
    setBoostersVisible(false)
    setBoostersLoading(true)
    ApiPostRequest({
      route: `/boosterpacks/open`,
      query: { tournament_id: props.params.tournamentID },
      body: {
        set_code: booster_pack.set_code
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
        booster_pack.available -= 1
        setLastOpenedPack(booster_pack)
      }
    })
  }

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

  return (
    <>
      {tournamentPlayer && <PackStoreModal
        tournamentID={props.params.tournamentID}
        tournamentPlayer={tournamentPlayer}
        availableCoins={tournamentPlayer?.game_resources.coins}
        isOpen={isPackStoreModalOpen}
        closeFn={() => setIsPackStoreModalOpen(false)}
        refreshBoostersFn={refreshData}
      />
      }
      <Layout tournamentID={props.params.tournamentID}>
        <div className="mx-16 my-16">
          {/* Open packs admin */}
          {/* TODO: Move this to an admin page */}
          {
            (tournamentPlayer?.access_level == "al_administrator" || tournamentPlayer?.access_level == "al_moderator") &&
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
          <Header title="Available packs" endContent={<Button size="md" color="success" onClick={() => setIsPackStoreModalOpen(true)}><BsFillCartFill />Buy packs</Button>} />
          {boostersLoading || !tournamentPlayer ? <div className="flex justify-center"> <Spinner /></div> :
            <div className="flex flex-row gap-2 justify-center">
              {boostersVisible && (
                <div className="bg-gray-800 w-[450px] border-small px-1 py-2 rounded-small border-default-200">
                  <Listbox>
                    {
                      tournamentPlayer.game_resources.booster_packs.map((booster_pack) =>
                        <ListboxItem
                          className="text-white"
                          onPress={() => sendOpenPackRequest(booster_pack)}
                          key={booster_pack.set_code}
                          startContent={<div className="text-gray-500 w-16">{booster_pack.set_code}</div>}
                          endContent={<div className="text-gray-500 text-right">{`(${booster_pack.available})`}</div>}
                        >
                          {booster_pack.name}
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
                  <ButtonGroup>
                    <Button
                      color="success"
                      isDisabled={lastOpenedPack && lastOpenedPack.available <= 0}
                      onClick={() => {
                        setBoostersLoading(true)
                        setCurrentCards([])
                        refreshData()
                        if (lastOpenedPack && lastOpenedPack.available > 0) {
                          sendOpenPackRequest(lastOpenedPack)
                        } else {
                          setBoostersVisible(true)
                        }
                      }}>
                      Open another ({lastOpenedPack?.available} {lastOpenedPack?.set_code} remaining)
                    </Button>
                    <Button
                      color="danger"
                      onClick={() => {
                        setCurrentCards([])
                        refreshData()
                        setBoostersVisible(true)
                      }}>
                      Close
                    </Button>
                  </ButtonGroup>
                </div>
              )}
            </div>
          }
        </div >
      </Layout >
    </>
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

type TournamentPlayersWithUsernames = {
  username: string
  tournamentPlayerId: string
}

function AddPacks(props: AddPacksProps) {
  let [availablePacks, setAvailablePacks] = useState<BoosterPack[]>([])
  let [error, setError] = useState<string>("")
  let [count, setCount] = useState<number>(0)
  let [setCode, setSetCode] = useState<string>("")
  let [availableTournamentPlayers, setAvailableTournamentPlayers] = useState<TournamentPlayersWithUsernames[]>([])
  let [tournamentPlayerId, setTournamentPlayerId] = useState<string>()
  let [loading, setLoading] = useState<boolean>(false)

  let isButtonDisabled = count <= 0 || setCode == ""

  let sendAddPacksRequest = () => {
    setError("")
    setLoading(true)
    ApiPostRequest({
      body: {
        count,
        set_code: setCode,
        tournament_player_id: tournamentPlayerId
      },
      route: "/boosterpacks/tournament",
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

  let refreshAvailableBoosters = () => {
    ApiGetRequest({
      route: "/boosterpacks/tournament",
      query: { tournament_id: props.tournamentID },
      responseHandler: (res: { booster_packs: BoosterPack[] }) => {
        setAvailablePacks(res.booster_packs)
      },
      errorHandler: (err) => {
        setError(err)
      }
    })
  }

  let getTournamentPlayers = () => {
    ApiGetRequest({
      route: "/tournament/tournament_player",
      query: { tournament_id: props.tournamentID },
      responseHandler: (res: { tournament_players: TournamentPlayer[], users: User[] }) => {
        let usersById: { [userId: string]: User } = {}
        res.users.forEach((user) => {
          usersById[user.id] = user
        })
        let tournamentPlayersWithUsernames = res.tournament_players.map((tournament_player) => {
          return {
            username: usersById[tournament_player.user_id].username,
            tournamentPlayerId: tournament_player.id,
          }
        })
        setAvailableTournamentPlayers(tournamentPlayersWithUsernames)
      },
      errorHandler: (err) => {
        setError(err)
      }
    })
  }

  useEffect(() => {
    refreshAvailableBoosters()
    getTournamentPlayers()
  }, [props.tournamentID])

  return (
    <div className="w-full">
      <form onSubmit={sendAddPacksRequest} className="flex flex-col gap-2 items-center w-full">
        <div className="flex flex-row gap-2 items-center w-full">
          <Input
            disabled={loading}
            onChange={(e) => setCount(Number(e.target.value))}
            variant="bordered"
            type="number"
            min={0}
            label="Amount"
            placeholder="Number of packs"
            labelPlacement="inside"
            className="text-white max-w-64"
            endContent={
              <div className="pointer-events-none flex items-center">
                <span className="text-gray-300 text-small">packs</span>
              </div>
            }
          />
          <Autocomplete
            disabled={loading}
            onInputChange={(value) => setSetCode(availablePacks.find(v => `${v.set_code} - ${v.name}` == value)?.set_code || "")}
            id="set"
            label="Booster pack"
            labelPlacement="inside"
            placeholder="Select a booster pack"
            className="text-white"
            defaultItems={availablePacks.map((val) => { return { value: val.set_code, label: `${val.set_code} - ${val.name}` } })}
          >
            {(item) => <AutocompleteItem className="text-white" key={item.value}>{item.label}</AutocompleteItem>}
          </Autocomplete>
          <Autocomplete
            disabled={loading}
            onInputChange={(value) => setTournamentPlayerId(availableTournamentPlayers.find(availablePlayer => availablePlayer.username === value)?.tournamentPlayerId)}
            id="tournamentPlayerId"
            label="Tournament Player"
            labelPlacement="inside"
            placeholder="All players"
            className="text-white max-w-xs"
            defaultItems={availableTournamentPlayers.map((val) => { return { value: val.tournamentPlayerId, label: `${val.username}` } })}
          >
            {(item) => <AutocompleteItem className="text-white" key={item.value}>{item.label}</AutocompleteItem>}
          </Autocomplete>
        </div>
        <Button isLoading={loading} onPress={sendAddPacksRequest} disabled={isButtonDisabled} size="md" className={`w-full ${isButtonDisabled || "bg-primary-500"}`} aria-label="Like">Add packs</Button>
      </form >
      <p className="text-sm font-light text-red-400">{error}</p>
    </div >
  )
}

type PackStoreModalProps = {
  tournamentID: string
  tournamentPlayer: TournamentPlayer
  availableCoins?: number
  isOpen: boolean
  closeFn: () => void
  refreshBoostersFn: () => void
}

function PackStoreModal(props: PackStoreModalProps) {
  let [error, setError] = useState<string>("")
  let [isLoading, setIsLoading] = useState<boolean>(false)
  let [store, setStore] = useState<Store>()
  let [availableBoosterPacks, setAvailableBoosterPacks] = useState<GlobalBoosterPack[]>([])

  let sendBuyBoosterRequest = (boosterPackId: string) => {
    setError("")
    setIsLoading(true)
    DoBuyBoosterPackRequest(
      boosterPackId,
      props.tournamentID,
      () => {
        setIsLoading(false)
        props.refreshBoostersFn()
      },
      (err) => {
        setIsLoading(false)
        setError(err)
        props.refreshBoostersFn()
      },
    )
  }

  let refreshData = () => {
    setIsLoading(true)
    setError("")
    DoGetTournamentStoreRequest(
      props.tournamentID,
      (s) => {
        setStore(s)
        setIsLoading(false)
      },
      (err) => {
        setError(err)
        setIsLoading(false)
      }
    )
    DoGetAvailableBoosterPacksRequest(
      props.tournamentID,
      (bp) => {
        setAvailableBoosterPacks(bp)
        setIsLoading(false)
      },
      (err) => {
        setError(err)
        setIsLoading(false)
      }
    )
  }

  useEffect(() => {
    refreshData()
  }, [props.isOpen])

  return (
    <Modal
      hideCloseButton
      onClose={props.closeFn}
      isOpen={props.isOpen}
      placement="top-center"
      size="2xl"
    >
      <ModalContent>
        <ModalHeader className="flex flex-col text-white gap-1">Booster pack store</ModalHeader>
        <ModalBody>
          {
            !store || isLoading ? <div className="flex flex-col"><Spinner /></div> : error ? error :
              <div className="bg-gray-800 border-small px-1 py-2 rounded-small border-default-200 w-full">
                <Listbox
                  className="w-full"
                  disabledKeys={
                    store.booster_packs.
                      filter(bp => bp.coin_price > props.tournamentPlayer.game_resources.coins).
                      map(bp => bp.booster_pack_id)
                  }
                >
                  {
                    store.booster_packs.map((booster_pack) =>
                      <ListboxItem
                        className="text-white"
                        onPress={() => sendBuyBoosterRequest(booster_pack.booster_pack_id)}
                        key={booster_pack.booster_pack_id}

                        startContent={<div className="text-gray-500 w-16">{availableBoosterPacks.filter(bp => bp.id == booster_pack.booster_pack_id)[0]?.set_code || ""}</div>}
                        endContent={<div className="text-orange-500 text-right flex flex-row items-center gap-2">{booster_pack.coin_price} <BsHypnotize className="text-orange-400" /></div>}
                      >
                        {availableBoosterPacks.filter(bp => bp.id == booster_pack.booster_pack_id)[0]?.name || ""}
                      </ListboxItem>
                    )
                  }

                </Listbox>
              </div>
          }
          <p className="text-sm font-light text-red-400 h-2">{error}</p>
        </ModalBody>
        <ModalFooter className="flex flex-row gap-4 justify-between items-center">
          <div className="flex flex-row items-center gap-2 text-white">
            <BsHypnotize className="text-orange-400" />
            {props.availableCoins} coins available
          </div>
          <Button isDisabled={isLoading} color="danger" variant="flat" onPress={props.closeFn}>
            Close
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}