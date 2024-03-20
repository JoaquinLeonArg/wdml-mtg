export type TournamentPlayer = {
  id: string
  tournament_id: string
  access_level: string
  user_id: string
  tournament_points: number
  game_resources: {
    booster_packs: BoosterPack[]
  }
}

export type BoosterPack = {
  available: number
  set_code: string
  name: string
  description: string
}

export type OwnedBoosterPack = {
  set_code: number
  count: number
}