import { DecklistCardProps } from "@/components/decklistcard"
import { OwnedCard } from "@/types/card"
import { DeckCard } from "@/types/deck"

export function groupCardsByType(deck_cards: DeckCard[], cardsById: { [id: string]: OwnedCard }) {
  let creatures: DecklistCardProps[] = []
  let artifacts: DecklistCardProps[] = []
  let enchantments: DecklistCardProps[] = []
  let sorceries: DecklistCardProps[] = []
  let instants: DecklistCardProps[] = []
  let lands: DecklistCardProps[] = []
  let other: DecklistCardProps[] = []

  deck_cards.forEach((card) => {
    if (cardsById[card.owned_card_id].card_data.types.includes("Creature")) {
      creatures.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.types.includes("Artifact")) {
      artifacts.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.types.includes("Enchantment")) {
      enchantments.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.types.includes("Sorcery")) {
      sorceries.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.types.includes("Instant")) {
      instants.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.types.includes("Land")) {
      lands.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else {
      other.push({ card: cardsById[card.owned_card_id], count: card.count })
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

export function groupCardsByColor(deck_cards: DeckCard[], cardsById: { [id: string]: OwnedCard }) {
  let white: DecklistCardProps[] = []
  let blue: DecklistCardProps[] = []
  let black: DecklistCardProps[] = []
  let red: DecklistCardProps[] = []
  let green: DecklistCardProps[] = []
  let multicolor: DecklistCardProps[] = []
  let colorless: DecklistCardProps[] = []
  let lands: DecklistCardProps[] = []

  deck_cards.forEach((card) => {
    if (cardsById[card.owned_card_id].card_data.colors.length > 1) {
      multicolor.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.colors.includes("W")) {
      white.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.colors.includes("U")) {
      blue.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.colors.includes("B")) {
      black.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.colors.includes("R")) {
      red.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.colors.includes("G")) {
      green.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.types.includes("Land")) {
      lands.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else {
      colorless.push({ card: cardsById[card.owned_card_id], count: card.count })
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


export function groupCardsByMV(deck_cards: DeckCard[], cardsById: { [id: string]: OwnedCard }) {
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

  deck_cards.forEach((card) => {
    if (cardsById[card.owned_card_id].card_data.mana_value == 0) {
      zero.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.mana_value == 1) {
      one.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.mana_value == 2) {
      two.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.mana_value == 3) {
      three.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.mana_value == 4) {
      four.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.mana_value == 5) {
      five.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.mana_value == 6) {
      six.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.mana_value == 7) {
      seven.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else if (cardsById[card.owned_card_id].card_data.mana_value == 8) {
      eight.push({ card: cardsById[card.owned_card_id], count: card.count })
    } else {
      nineplus.push({ card: cardsById[card.owned_card_id], count: card.count })
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
    "9+ mana": nineplus
  }
}