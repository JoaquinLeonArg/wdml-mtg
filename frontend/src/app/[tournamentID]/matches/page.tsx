"use client"

import { Header } from "@/components/header"
import Layout from "@/components/layout"
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests"
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
      responseHandler: (res: { tournament_players: TournamentPlayer[] }) => {
        let tournamentPlayersByID: { [tournament_player_id: string]: TournamentPlayer } = {}
        res.tournament_players.forEach((tp) => tournamentPlayersByID[tp.id] = tp)
        setTournamentPlayersByID(tournamentPlayersByID)
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
        <Header title="Players" />
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
                          <>tournamentPlayersByID[pd.player_id]</>
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
    </Layout>
  )
}
