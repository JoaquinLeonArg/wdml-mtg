export type Tournament = {
  id: string
  invite_code: string
  name: string
  description: string
  store: Store
}

export type Store = {
  booster_packs: StoreBoosterPack[]
}

export type StoreBoosterPack = {
  booster_pack_id: string
  coin_price: number
}
