export type Deck = {
  id: string
  tournament_player_id: string
  name: string
  description: string
  cards: DeckCard[]
}

export type DeckCard = {
  owned_card_id: string
  count: number
  board: string
}