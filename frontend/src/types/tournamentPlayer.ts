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
  data: BoosterPackData
}

export type BoosterPackData = {
  set_name: string
  set_code: string
}