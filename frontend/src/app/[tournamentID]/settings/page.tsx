"use client"

import { Header, MiniHeader } from "@/components/header"
import Layout from "@/components/layout"
import { DoCreateBoosterPackRequest, DoGetAvailableBoosterPacksRequest, DoUpdateBoosterPackRequest } from "@/requests/boosterpacks"
import { DoGetTournamentStoreRequest, DoUpdateTournamentStoreRequest } from "@/requests/tournament"
import { BoosterPack } from "@/types/boosterPack"
import { Store, StoreBoosterPack } from "@/types/tournament"
import { Autocomplete, AutocompleteItem, Button, Input, Spinner, Textarea } from "@nextui-org/react"
import { useEffect, useState } from "react"

const DEFAULT_BOOSTER_PACK_JSON = `{
  "card_count": 15,
  "description": "Strixhaven Draft Booster Pack",
  "name": "Strixhaven, School of Mages",
  "set_code": "stx",
  "slots": [
    {
      "filter": "set:stx rarity:c -type:basic",
      "count": 8
    },
    {
      "options": [
        { "filter": "rarity:c -type:basic", "weight": 66 },
        { "filter": "rarity:c", "weight": 20 },
        { "filter": "rarity:u", "weight": 10 },
        { "filter": "rarity:r", "weight": 3 },
        { "filter": "rarity:m", "weight": 1 }
      ],
      "filter": "set:stx",
      "count": 1
    },
    {
      "filter": "set:stx rarity:u",
      "count": 3
    },
    {
      "options": [
        { "filter": "rarity:c", "weight": 30 },
        { "filter": "rarity:r", "weight": 7 },
        { "filter": "rarity:m", "weight": 1 }
      ],
      "filter": "set:stx type:Lesson",
      "count": 1
    },
    {
      "options": [
        { "filter": "rarity:r", "weight": 7 },
        { "filter": "rarity:m", "weight": 1 }
      ],
      "filter": "set:stx",
      "count": 1
    },
    {
      "options": [
        { "filter": "rarity:u", "weight": 30 },
        { "filter": "rarity:r", "weight": 7 },
        { "filter": "rarity:m", "weight": 1 }
      ],
      "filter": "set:sta",
      "count": 1
    }
  ]
}`

export default function ConfigPage(props: any) {
  let [availableBoosterPacks, setAvailableBoosterPacks] = useState<BoosterPack[]>([])
  let [store, setStore] = useState<Store>()
  let [isLoadingAvailable, setIsLoadingAvailable] = useState<boolean>(true)
  let [isLoadingStore, setIsLoadingStore] = useState<boolean>(true)
  let [error, setError] = useState<string>("")
  let [boosterPackJson, setBoosterPackJson] = useState<string>(DEFAULT_BOOSTER_PACK_JSON)
  let [boosterPackJsonError, setBoosterPackJsonError] = useState<string>("")
  let [isLoadingBoosterPack, setIsLoadingBoosterPack] = useState<boolean>(false)

  let refreshData = () => {
    setError("")
    setIsLoadingAvailable(true)
    setIsLoadingStore(true)
    DoGetAvailableBoosterPacksRequest(
      props.params.tournamentID,
      (booster_packs) => {
        setAvailableBoosterPacks(booster_packs)
        setIsLoadingAvailable(false)
      },
      (err) => {
        setError(err)
        setIsLoadingAvailable(false)
      },
    )
    DoGetTournamentStoreRequest(
      props.params.tournamentID,
      (store) => {
        setStore(store)
        setIsLoadingStore(false)
      },
      (err) => {
        setError(err)
        setIsLoadingStore(false)
      }
    )
  }

  let sendUpdateStoreRequest = () => {
    if (!store) return
    setError("")
    setIsLoadingAvailable(true)
    setIsLoadingStore(true)
    DoUpdateTournamentStoreRequest(
      props.params.tournamentID,
      store,
      () => { refreshData() },
      (err) => { setError(err) }
    )
  }

  let validateBoosterPackJson = () => {
    setBoosterPackJsonError("")
    try {
      JSON.parse(boosterPackJson)
    } catch (e) {
      setBoosterPackJsonError("Invalid JSON: " + e)
      return false
    }
    return true
  }

  let sendCreateBoosterPack = () => {
    if (!validateBoosterPackJson()) return
    setIsLoadingBoosterPack(true)
    DoCreateBoosterPackRequest(
      JSON.parse(boosterPackJson),
      () => {
        setIsLoadingBoosterPack(false)
        setBoosterPackJson("")
        refreshData()
      },
      (err) => {
        setIsLoadingBoosterPack(false)
        switch (err) {
          case "DUPLICATED_RESOURCE":
            setBoosterPackJsonError("A booster pack with this set_code already exists. Use Update instead.")
            break
          default:
            setBoosterPackJsonError(err)
        }
      }
    )
  }

  let sendUpdateBoosterPack = () => {
    if (!validateBoosterPackJson()) return
    setIsLoadingBoosterPack(true)
    DoUpdateBoosterPackRequest(
      JSON.parse(boosterPackJson),
      () => {
        setIsLoadingBoosterPack(false)
        setBoosterPackJson("")
        refreshData()
      },
      (err) => {
        setIsLoadingBoosterPack(false)
        setBoosterPackJsonError(err)
      }
    )
  }

  useEffect(() => {
    refreshData()
  }, [props.params.tournamentID])


  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        {
          error ? error : isLoadingStore || isLoadingAvailable ? <div className="flex flex-col"><Spinner /></div> :
            <>
              <Header title="Settings" />
              <MiniHeader title="Store" />
              <div className="flex flex-col gap-4 mb-4">
                {store != undefined && availableBoosterPacks &&
                  store.booster_packs.map((booster_pack: StoreBoosterPack, index: number) => {
                    let boosterPackData = availableBoosterPacks.filter(bp => bp.id == booster_pack.booster_pack_id)[0]
                    return (
                      <div key={index} className="flex flex-row gap-2 items-center w-full">
                        <Input
                          onChange={(e) => {
                            if (!store) return
                            let newStoreBoosterPacks = [...store.booster_packs]
                            newStoreBoosterPacks[index].coin_price = Number(e.target.value)
                            setStore({ ...store, booster_packs: newStoreBoosterPacks })
                          }}
                          variant="bordered"
                          type="number"
                          min={0}
                          label="Cost"
                          value={String(booster_pack.coin_price)}
                          placeholder="Pack cost"
                          labelPlacement="inside"
                          className="text-white max-w-64"
                          endContent={
                            <div className="pointer-events-none flex items-center">
                              <span className="text-gray-300 text-small">coins</span>
                            </div>
                          }
                        />
                        <Autocomplete
                          onInputChange={(value) => {
                            if (!store) return
                            let newStoreBoosterPacks = [...store.booster_packs]
                            newStoreBoosterPacks[index].booster_pack_id = availableBoosterPacks.find(v => `${v.set_code} - ${v.name}` == value)?.id || ""
                            setStore({ ...store, booster_packs: newStoreBoosterPacks })
                          }}
                          id="set"
                          label="Booster type"
                          labelPlacement="inside"
                          placeholder="Select a booster pack"
                          className="text-white"
                          defaultSelectedKey={boosterPackData?.set_code || ""}
                          defaultItems={availableBoosterPacks.map((val) => { return { value: val.set_code, label: `${val.set_code} - ${val.name}` } })}
                        >
                          {(item) => <AutocompleteItem className="text-white" key={item.value}>{item.label}</AutocompleteItem>}
                        </Autocomplete>
                        <Button
                          onClick={() => {
                            if (!store) return
                            let newStoreBoosterPacks = [...store.booster_packs]
                            newStoreBoosterPacks.splice(index, 1)
                            setStore({ ...store, booster_packs: newStoreBoosterPacks })
                          }}
                          color="danger"
                          isIconOnly
                        >
                          X
                        </Button>
                      </div>
                    )
                  })}
                {
                  store != undefined &&
                  <Button
                    onPress={() => {
                      if (!store) return
                      setStore({ ...store, booster_packs: [...store.booster_packs, { booster_pack_id: "", coin_price: 0 }] })
                    }
                    }
                    size="md" color="success" aria-label="Update"
                  >
                    Add store item
                  </Button>
                }
              </div>
              <Button onPress={sendUpdateStoreRequest} size="md" color="success" aria-label="Update">Update</Button>
              <MiniHeader title="Add/Edit Booster Packs" />
              <div className="flex flex-col gap-2 mb-4">
                <Textarea
                  label="Booster pack definition"
                  minRows={20}
                  maxRows={40}
                  value={boosterPackJson}
                  onValueChange={setBoosterPackJson}
                  className="text-white"
                />
                <p className="text-sm font-light text-red-400">{boosterPackJsonError}</p>
                <div className="flex flex-row gap-2">
                  <Button isLoading={isLoadingBoosterPack} onPress={sendCreateBoosterPack} color="success">Create</Button>
                  <Button isLoading={isLoadingBoosterPack} onPress={sendUpdateBoosterPack} color="warning">Update</Button>
                </div>
              </div>
            </>
        }
      </div>
    </Layout >
  )
}
