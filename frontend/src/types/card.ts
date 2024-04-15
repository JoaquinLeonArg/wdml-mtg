export type CardData = {
  rarity: "common" | "uncommon" | "rare" | "mythic" | "special"
  types: string[]
  colors: string[]
  mana_value: number
  mana_cost: string
  name: string
  image_url: string
  back_image_url: string
}

export type OwnedCard = {
  id: string
  count: number
  card_data: CardData
}

