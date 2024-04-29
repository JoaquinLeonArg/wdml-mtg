import { Store } from "@/types/tournament";
import { ApiGetRequest, ApiPostRequest } from "./requests";

export function DoGetTournamentStoreRequest(
  tournamentId: string,
  onResponse: (store: Store) => void,
  onError: (_: string) => void
) {
  ApiGetRequest({
    route: "/tournament/store",
    query: {
      tournament_id: tournamentId
    },
    responseHandler: (res: { store: Store }) => onResponse(res.store),
    errorHandler: (err: string) => onError(err)
  })
}

export function DoUpdateTournamentStoreRequest(
  tournamentId: string,
  store: Store,
  onResponse: () => void,
  onError: (_: string) => void
) {
  ApiPostRequest({
    route: "/tournament/store/update",
    query: {
      tournament_id: tournamentId
    },
    body: {
      store
    },
    responseHandler: () => onResponse(),
    errorHandler: (err: string) => onError(err)
  })
}