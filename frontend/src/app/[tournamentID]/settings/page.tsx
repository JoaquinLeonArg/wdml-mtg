"use client"

import { Header, MiniHeader } from "@/components/header"
import Layout from "@/components/layout"
import { DoGetAvailableBoosterPacksRequest } from "@/requests/boosterpacks"
import { DoGetTournamentStoreRequest, DoUpdateTournamentStoreRequest } from "@/requests/tournament"
import { BoosterPack } from "@/types/boosterPack"
import { Store, StoreBoosterPack } from "@/types/tournament"
import { Autocomplete, AutocompleteItem, Button, Input, Spinner } from "@nextui-org/react"
import { useEffect, useState } from "react"

export default function ConfigPage(props: any) {
  let [availableBoosterPacks, setAvailableBoosterPacks] = useState<BoosterPack[]>([])
  let [store, setStore] = useState<Store>()
  let [isLoadingAvailable, setIsLoadingAvailable] = useState<boolean>(true)
  let [isLoadingStore, setIsLoadingStore] = useState<boolean>(true)
  let [error, setError] = useState<string>("")

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

  useEffect(() => {
    refreshData()
  }, [props.params.tournamentID])

  console.log(store)

  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        {
          error ? error : isLoadingStore || isLoadingAvailable ? <div className="flex flex-col"><Spinner /></div> :
            <>
              <Header title="Settings" />
              <MiniHeader title="Store" />
              <div className="flex flex-col gap-4 mb-4">
                {store && availableBoosterPacks &&
                  store.booster_packs.map((booster_pack: StoreBoosterPack, index: number) => {
                    let boosterPackData = availableBoosterPacks.filter(bp => bp.id == booster_pack.booster_pack_id)[0]
                    return (
                      <div className="flex flex-row gap-2 items-center w-full">
                        <Input
                          onChange={(e) => {
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
                  store &&
                  <Button
                    onPress={() => setStore({ ...store, booster_packs: [...store.booster_packs, { booster_pack_id: "", coin_price: 0 }] })}
                    size="md" color="success" aria-label="Update"
                  >
                    Add store item
                  </Button>
                }
              </div>
              <Button onPress={sendUpdateStoreRequest} size="md" color="success" aria-label="Update">Update</Button>
            </>
        }
      </div>
    </Layout >
  )
}
