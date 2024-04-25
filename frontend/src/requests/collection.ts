import { CardData } from "@/types/card";
import { ApiPostRequest } from "./requests";

export function DoTradeUpRequest(
  cards: { [card_id: string]: number },
  tournamentId: string,
  onResponse: (_: CardData[]) => void,
  onError: (_: string) => void
) {
  ApiPostRequest({
    route: "/collection/tradeup",
    query: {
      tournament_id: tournamentId
    },
    body: {
      "cards": cards,
    },
    responseHandler: (resp: { cards: CardData[] }) => onResponse(resp.cards),
    errorHandler: (err: string) => onError(err)
  })
}