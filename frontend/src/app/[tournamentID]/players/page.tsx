"use client"

import { Header } from "@/components/header"
import Layout from "@/components/layout"
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests"
import { Deck } from "@/types/deck"
import { TournamentPlayer } from "@/types/tournamentPlayer"
import { User } from "@/types/user"
import { Button, Input, Listbox, ListboxItem, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, Spinner, Table, TableBody, TableCell, TableColumn, TableHeader, TableRow, Textarea, getKeyValue } from "@nextui-org/react"
import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"



export default function PlayersPage(props: any) {
  let [rows, setRows] = useState<{ username: string, tournament_points: number, coins: number }[]>([])
  let [isLoading, setIsLoading] = useState<boolean>(true)
  let [error, setError] = useState<string>("")

  let refreshData = () => {
    setIsLoading(true)
    ApiGetRequest({
      route: "/tournament/tournament_player",
      query: { tournament_id: props.params.tournamentID },
      errorHandler: (err) => {
        setError(err)
        setIsLoading(false)
      },
      responseHandler: (res: { tournament_players: TournamentPlayer[], users: User[] }) => {
        let usersById: { [userId: string]: User } = {}
        res.users.forEach((user) => {
          usersById[user.id] = user
        })
        let r = res.tournament_players.map((tournament_player) => {
          return {
            username: usersById[tournament_player.user_id].username,
            tournament_points: tournament_player.tournament_points,
            coins: tournament_player.game_resources.coins,
          }
        })
        setRows(r)
        setIsLoading(false)
      }
    })
  }

  useEffect(() => {
    refreshData()
  }, [props.params.tournamentID])

  let columns = [
    {
      key: "username",
      label: "Username"
    },
    {
      key: "tournament_points",
      label: "Points"
    },
    {
      key: "coins",
      label: "Coins"
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
            <TableBody items={rows}>
              {(item) => (
                <TableRow key={item.username}>
                  {(columnKey) => <TableCell className="text-white">{getKeyValue(item, columnKey)}</TableCell>}
                </TableRow>
              )}
            </TableBody>
          </Table>
        }
      </div>
    </Layout>
  )
}
