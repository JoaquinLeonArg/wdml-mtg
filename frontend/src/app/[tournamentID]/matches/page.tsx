"use client"

import { Header } from "@/components/header"
import Layout from "@/components/layout"
import { ApiGetRequest, ApiPostRequest, ApiPutRequest } from "@/requests/requests"
import { Deck } from "@/types/deck"
import { Match } from "@/types/match"
import { Season } from "@/types/season"
import { TournamentPlayer } from "@/types/tournamentPlayer"
import { User } from "@/types/user"
import { Button, Input, Listbox, ListboxItem, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, Spinner, Table, TableBody, TableCell, TableColumn, TableHeader, TableRow, Textarea, getKeyValue } from "@nextui-org/react"
import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"



export default function MatchesPage(props: any) {
  let [isLoading, setIsLoading] = useState<boolean>(true)
  let [error, setError] = useState<string>("")
  let [matches, setMatches] = useState<Match[]>([])
  let [tournamentPlayersByID, setTournamentPlayersByID] = useState<{ [tournament_player_id: string]: TournamentPlayer }>({})
  let [usersById, setUsersByID] = useState<{ [user_id: string]: User }>({})
  let [tournamentPlayer, setTournamentPlayer] = useState<TournamentPlayer>()
  let [season, setSeason] = useState<Season>()

  let refreshData = () => {
    setIsLoading(true)
    ApiGetRequest({
      route: "/season/all",
      query: { tournament_id: props.params.tournamentID },
      errorHandler: (err) => {
        setError(err)
        setIsLoading(false)
      },
      responseHandler: (res: { seasons: Season[] }) => {
        if (!res.seasons) { setIsLoading(false); return }
        let lastSeason = res.seasons.at(-1)
        if (!lastSeason) { setIsLoading(false); return }
        setSeason(lastSeason)
        ApiGetRequest({
          route: "/match",
          query: { season_id: lastSeason.id },
          errorHandler: (err) => {
            setError(err)
            setIsLoading(false)
          },
          responseHandler: (res: { matches: Match[] }) => {
            setMatches(res.matches)
          }
        })
        setIsLoading(false)
      }
    })
    ApiGetRequest({
      route: "/tournament/tournament_player",
      query: { tournament_id: props.params.tournamentID },
      errorHandler: (err) => {
      },
      responseHandler: (res: { tournament_players: TournamentPlayer[], users: User[] }) => {
        let tournamentPlayersByID: { [tournament_player_id: string]: TournamentPlayer } = {}
        res.tournament_players.forEach((tp) => tournamentPlayersByID[tp.id] = tp)
        setTournamentPlayersByID(tournamentPlayersByID)

        let usersByID: { [user_id: string]: User } = {}
        res.users.forEach((u) => usersByID[u.id] = u)
        setUsersByID(usersByID)
      }
    })
    ApiGetRequest({
      route: "/seasons/match",
      query: { season_id: props.params.tournamentID },
      errorHandler: (err) => {
      },
      responseHandler: (res: { matches: Match[] }) => {
        setMatches(res.matches)
      }
    })
    ApiGetRequest({
      route: "/tournament_player/tournament",
      query: { tournament_id: props.params.tournamentID },
      errorHandler: (err) => {
        setIsLoading(false)
      },
      responseHandler: (res: { tournament_player: TournamentPlayer }) => {
        setIsLoading(false)
        setTournamentPlayer(res.tournament_player)
      }
    })
  }

  let updateMatch = (matchID: string, tournamentPlayerID: string, pointsChange: number) => {
    let match = matches.find((m) => m.id == matchID)
    if (!match) return
    let players_data: { [tournament_player_id: string]: number } = {}
    match.players_data.forEach((pd) => players_data[pd.tournament_player_id] = pd.wins)
    players_data[tournamentPlayerID] += pointsChange

    setIsLoading(true)

    ApiPutRequest({
      route: "/match",
      query: { season_id: season?.id, match_id: matchID },
      body: {
        players_points: players_data,
        games_played: match.games_played + pointsChange
      },
      errorHandler: (err) => {
        setIsLoading(false)
      },
      responseHandler: () => {
        setIsLoading(false)
        refreshData()
      }
    })
  }

  useEffect(() => {
    refreshData()
  }, [props.params.tournamentID])

  let columns = [
    {
      key: "players",
      label: "Players"
    },
    {
      key: "wins",
      label: "Wins"
    }
  ]

  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        <Header title="Matches" />
        {isLoading ? <div className="flex justify-center"> <Spinner /></div> :

          <Table aria-label="Example table with dynamic content">
            <TableHeader columns={columns}>
              {(column) => <TableColumn key={column.key}>{column.label}</TableColumn>}
            </TableHeader>
            <TableBody items={matches || []}>
              {(item) => (
                <TableRow key={item.id}>
                  <TableCell className="text-white">
                    <div className="flex flex-col gap-2">
                      {item.players_data.map((pd) => {
                        return (
                          <div key={pd.tournament_player_id + item.id}>{usersById[tournamentPlayersByID[pd.tournament_player_id].user_id].username}</div>
                        )
                      })}
                    </div>
                  </TableCell>
                  <TableCell className="text-white">
                    <div className="flex flex-col gap-2">
                      {item.players_data.map((pd) => {
                        return (
                          tournamentPlayer && (
                            item.players_data.map((pd) => pd.tournament_player_id).includes(tournamentPlayer.id) ||
                            tournamentPlayer?.access_level == "al_administrator" || tournamentPlayer.access_level == "al_moderator") ?
                            <div className="flex flex-row gap-2 items-center">
                              <Button onClick={() => updateMatch(item.id, pd.tournament_player_id, -1)} size="sm" isIconOnly>-</Button>
                              <div> {pd.wins}</div>
                              <Button onClick={() => updateMatch(item.id, pd.tournament_player_id, 1)} size="sm" isIconOnly>+</Button>
                            </div> :
                            <div> {pd.wins}</div>
                        )
                      })}
                    </div>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        }
      </div>
    </Layout >
  )
}
