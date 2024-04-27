import { DecklistCardListProps, DecklistCardProps } from "@/components/decklistcard"
import { OwnedCard } from "@/types/card"
import { DeckCard } from "@/types/deck"

export function groupCardsByType(deckCards: DeckCard[], allDeckCards: DeckCard[], cardsById: { [id: string]: OwnedCard }) {
  let creatures: DecklistCardProps[] = []
  let artifacts: DecklistCardProps[] = []
  let enchantments: DecklistCardProps[] = []
  let sorceries: DecklistCardProps[] = []
  let instants: DecklistCardProps[] = []
  let lands: DecklistCardProps[] = []
  let other: DecklistCardProps[] = []

  deckCards.forEach((card) => {
    let types = cardsById[card.owned_card_id].card_data.types
    let totalCount = allDeckCards.reduce((prev, deckCard) => card.owned_card_id == deckCard.owned_card_id ? prev + deckCard.count : prev, 0)
    console.log(totalCount)
    if (types.includes("Land")) {
      lands.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (types.includes("Creature")) {
      creatures.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (types.includes("Artifact")) {
      artifacts.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (types.includes("Enchantment")) {
      enchantments.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (types.includes("Sorcery")) {
      sorceries.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (types.includes("Instant")) {
      instants.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else {
      other.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    }
  })

  return {
    "Creatures": creatures,
    "Artifacts": artifacts,
    "Enchantments": enchantments,
    "Sorceries": sorceries,
    "Instants": instants,
    "Lands": lands,
    "Other": other
  }
}

export function groupCardsByColor(deckCards: DeckCard[], allDeckCards: DeckCard[], cardsById: { [id: string]: OwnedCard }) {
  let white: DecklistCardProps[] = []
  let blue: DecklistCardProps[] = []
  let black: DecklistCardProps[] = []
  let red: DecklistCardProps[] = []
  let green: DecklistCardProps[] = []
  let multicolor: DecklistCardProps[] = []
  let colorless: DecklistCardProps[] = []
  let lands: DecklistCardProps[] = []

  deckCards.forEach((card) => {
    let colors = cardsById[card.owned_card_id].card_data.colors
    let totalCount = allDeckCards.reduce((prev, deckCard) => card.owned_card_id == deckCard.owned_card_id ? prev + deckCard.count : prev, 0)
    if (colors.length > 1) {
      multicolor.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (colors.includes("W")) {
      white.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (colors.includes("U")) {
      blue.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (colors.includes("B")) {
      black.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (colors.includes("R")) {
      red.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (colors.includes("G")) {
      green.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (cardsById[card.owned_card_id].card_data.types.includes("Land")) {
      lands.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else {
      colorless.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    }
  })

  return {
    "White": white,
    "Blue": blue,
    "Black": black,
    "Red": red,
    "Green": green,
    "Multicolor": multicolor,
    "Colorless": colorless,
    "Lands": lands
  }
}


export function groupCardsByMV(deckCards: DeckCard[], allDeckCards: DeckCard[], cardsById: { [id: string]: OwnedCard }) {
  let land: DecklistCardListProps[] = []
  let zero: DecklistCardProps[] = []
  let one: DecklistCardProps[] = []
  let two: DecklistCardProps[] = []
  let three: DecklistCardProps[] = []
  let four: DecklistCardProps[] = []
  let five: DecklistCardProps[] = []
  let six: DecklistCardProps[] = []
  let seven: DecklistCardProps[] = []
  let eight: DecklistCardProps[] = []
  let nineplus: DecklistCardProps[] = []

  console.log(deckCards)
  deckCards.forEach((card) => {
    let mv = cardsById[card.owned_card_id].card_data.mana_value
    let totalCount = allDeckCards.reduce((prev, deckCard) => card.owned_card_id == deckCard.owned_card_id ? prev + deckCard.count : prev, 0)
    if (cardsById[card.owned_card_id].card_data.types.includes("Land")) {
      land.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (mv == 0) {
      zero.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (mv == 1) {
      one.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (mv == 2) {
      two.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (mv == 3) {
      three.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (mv == 4) {
      four.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (mv == 5) {
      five.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (mv == 6) {
      six.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (mv == 7) {
      seven.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else if (mv == 8) {
      eight.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    } else {
      nineplus.push({ card: cardsById[card.owned_card_id], count: card.count, totalCount })
    }
  })

  return {
    "1 mana": one,
    "2 mana": two,
    "3 mana": three,
    "4 mana": four,
    "5 mana": five,
    "6 mana": six,
    "7 mana": seven,
    "8 mana": eight,
    "9+ mana": nineplus,
    "Lands": land
  }
}