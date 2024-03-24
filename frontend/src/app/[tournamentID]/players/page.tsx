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
import { BsFillSendFill } from "react-icons/bs"



export default function PlayersPage(props: any) {
  let [originalRows, setOriginalRows] = useState<{ id: string, username: string, tournament_points: number, coins: number }[]>([])
  let [rows, setRows] = useState<{ id: string, username: string, tournament_points: number, coins: number }[]>([])
  let [isLoading, setIsLoading] = useState<boolean>(true)
  let [error, setError] = useState<string>("")
  let [tournamentPlayer, setTournamentPlayer] = useState<TournamentPlayer>()

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
            id: tournament_player.id,
            username: usersById[tournament_player.user_id].username,
            tournament_points: tournament_player.tournament_points,
            coins: tournament_player.game_resources.coins,
          }
        })
        setRows([...r])
        setOriginalRows([...r])
        setIsLoading(false)
      }
    })
    ApiGetRequest({
      route: "/tournament_player/tournament",
      query: { tournament_id: props.params.tournamentID },
      errorHandler: (err) => {
      },
      responseHandler: (res: { tournament_player: TournamentPlayer }) => {
        setTournamentPlayer(res.tournament_player)
      }
    })
  }

  let sendChangeTournamentPointsRequest = (user_id: string) => {
    let tournamentPlayerIndex = rows.findIndex((user) => user.id == user_id)
    if (tournamentPlayerIndex == -1) { return } // TODO: Show error
    setIsLoading(true)
    ApiPostRequest({
      route: "/tournament_player/points",
      query: { tournament_id: props.params.tournamentID, tournament_player_id: originalRows[tournamentPlayerIndex].id },
      body: {
        points: rows[tournamentPlayerIndex].tournament_points - originalRows[tournamentPlayerIndex].tournament_points
      },
      errorHandler: () => {
        setIsLoading(false)
      },
      responseHandler: () => {
        refreshData()
        setIsLoading(false)
      }
    })
  }

  let sendChangeCoinsRequest = (user_id: string) => {
    let tournamentPlayerIndex = rows.findIndex((user) => user.id == user_id)
    if (tournamentPlayerIndex == -1) { return } // TODO: Show error
    setIsLoading(true)
    ApiPostRequest({
      route: "/tournament_player/coins",
      query: { tournament_id: props.params.tournamentID, tournament_player_id: originalRows[tournamentPlayerIndex].id },
      body: {
        coins: rows[tournamentPlayerIndex].coins - originalRows[tournamentPlayerIndex].coins
      },
      errorHandler: () => {
        setIsLoading(false)
      },
      responseHandler: () => {
        refreshData()
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
                  {(columnKey) => {
                    if (columnKey == "tournament_points") {
                      return <TableCell className="text-white">
                        <div className="flex flex-row items-center gap-2">
                          {
                            (tournamentPlayer?.access_level == "al_administrator" || tournamentPlayer?.access_level == "al_moderator") ?
                              <>
                                <Input
                                  defaultValue={getKeyValue(item, columnKey)}
                                  disabled={isLoading}
                                  onChange={(e) => {
                                    let newRows = [...rows]
                                    let tournamentPlayerIndex = newRows
                                      .findIndex((r) => r.id == item.id)
                                    let newRow = { ...newRows[tournamentPlayerIndex] }
                                    newRow.tournament_points = Number(e.target.value)
                                    newRows[tournamentPlayerIndex] = newRow
                                    setRows(newRows)
                                  }}
                                  variant="bordered"
                                  size="sm"
                                  type="number"
                                  placeholder="0"
                                  className="text-white max-w-16"
                                />
                                <Button size="sm" color="success" isIconOnly onClick={() => {
                                  sendChangeTournamentPointsRequest(item.id)
                                }}><BsFillSendFill /></Button>
                              </> :
                              <div>{getKeyValue(item, columnKey)}</div>
                          }
                        </div>
                      </TableCell>
                    }
                    if (columnKey == "coins") {
                      return <TableCell className="text-white">
                        <div className="flex flex-row items-center gap-2">
                          {
                            tournamentPlayer?.access_level == "al_administrator" || tournamentPlayer?.access_level == "al_moderator" ?
                              <>
                                <Input
                                  defaultValue={getKeyValue(item, columnKey)}
                                  disabled={isLoading}
                                  onChange={(e) => {
                                    let newRows = [...rows]
                                    let tournamentPlayerIndex = newRows
                                      .findIndex((r) => r.id == item.id)
                                    let newRow = { ...newRows[tournamentPlayerIndex] }
                                    newRow.coins = Number(e.target.value)
                                    newRows[tournamentPlayerIndex] = newRow
                                    setRows(newRows)
                                  }}
                                  variant="bordered"
                                  size="sm"
                                  type="number"
                                  placeholder="0"
                                  className="text-white max-w-16"
                                />
                                <Button size="sm" color="success" isIconOnly onClick={() => {
                                  sendChangeCoinsRequest(item.id)
                                }}><BsFillSendFill /></Button>
                              </> :
                              <div>{getKeyValue(item, columnKey)}</div>
                          }
                        </div>
                      </TableCell>
                    }
                    return (<TableCell className="text-white">{getKeyValue(item, columnKey)}</TableCell>)
                  }
                  }
                </TableRow>
              )}
            </TableBody>
          </Table>
        }
      </div>
    </Layout>
  )
}
