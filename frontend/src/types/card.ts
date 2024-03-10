export type CardData = {
  card_rarity: string
  image_url: string
}

export type OwnedCard = {
  count: number
  card_data: CardData
}