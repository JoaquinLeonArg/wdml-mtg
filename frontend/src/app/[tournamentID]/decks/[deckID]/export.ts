import { OwnedCard } from "@/types/card";
import { Deck } from "@/types/deck";

export function FormatDeck(deck: Deck, cards: OwnedCard[]) {
  let cardsById: { [id: string]: OwnedCard } = {}
  cards.forEach((card) => {
    cardsById[card.id] = card
  })

  let text = ""
  deck.cards
    .filter((card) => card.board == "b_mainboard")
    .forEach(
      (card) => {
        text += `${card.count} ${cardsById[card.owned_card_id].card_data.name}\n`
      }
    )
  text += "\n" // Sideboard separator for Cockatrice
  deck.cards
    .filter((card) => card.board == "b_sideboard")
    .forEach(
      (card) => {
        text += `${card.count} ${cardsById[card.owned_card_id].card_data.name}\n`
      }
    )

  return text
}