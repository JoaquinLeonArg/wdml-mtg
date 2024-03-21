export type TournamentPost = {
  id: string
  tournament_id: string
  tournament_player_id: string
  title: string
  blocks: TournamentPostBlock[]

}

export type TournamentPostBlock = {
  title: string
  collapsable: boolean
  content: TournamentPostContent[]
}

export type TournamentPostContent = {
  type: "tbct_text" | "tbct_image"
  content: string
  extra_data: { [key: string]: string }
}