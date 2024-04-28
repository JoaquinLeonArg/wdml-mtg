import { BoosterPack } from "@/types/boosterPack";
import { ApiGetRequest, ApiPostRequest } from "./requests";

export function DoBuyBoosterPackRequest(
  boosterPackId: string,
  tournamentId: string,
  onResponse: () => void,
  onError: (_: string) => void
) {
  ApiPostRequest({
    route: "/boosterpacks/buy",
    query: {
      tournament_id: tournamentId
    },
    body: {
      "booster_pack_id": boosterPackId,
    },
    responseHandler: (_) => onResponse(),
    errorHandler: (err: string) => onError(err)
  })
}

export function DoGetAvailableBoosterPacksRequest(
  tournamentId: string,
  onResponse: (booster_packs: BoosterPack[]) => void,
  onError: (_: string) => void
) {
  ApiGetRequest({
    route: "/boosterpacks/tournament",
    query: {
      tournament_id: tournamentId
    },
    responseHandler: (res: { booster_packs: BoosterPack[] }) => onResponse(res.booster_packs),
    errorHandler: (err: string) => onError(err)
  })
}