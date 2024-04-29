export type OwnedBoosterPack = {
  available: number
  data: {
    set_code: string
    set_name: string
    booster_type: string
    description: string
  }
}

export type BoosterPack = {
  id: string
  set_code: string
  name: string
  description: string
  card_count: number
}