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
  set_name: string
  set_code: string
  count: number
  data: BoosterPackData
}

export type BoosterPackData = {

}