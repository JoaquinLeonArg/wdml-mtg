export type Match = {
  id: string
  season_id: string
  players_data: MatchPlayerData[]
  games_played: number
  gamemode: "gm_1v1" | "gm_arc" | "gm_hop" | "gm_edh" | "gm_2hg"
  completed: boolean
}

export type MatchPlayerData = {
  tournament_player_id: string
  wins: number
  tags: string[]
}