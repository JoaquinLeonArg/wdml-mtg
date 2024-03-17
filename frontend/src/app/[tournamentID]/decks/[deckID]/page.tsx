"use client"

import { Header, MiniHeader } from "@/components/header"
import Layout from "@/components/layout"
import { ApiGetRequest } from "@/requests/requests"
import { Deck } from "@/types/deck"
import { Button, Checkbox, Input, Link, Listbox, ListboxItem, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, useDisclosure } from "@nextui-org/react"
import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"



export default function EditDeckPage(props: any) {
  let router = useRouter()
  let [deck, setDeck] = useState<Deck>()
  let [deckLoading, setDeckLoading] = useState<boolean>(false)

  useEffect(() => refreshData())

  let refreshData = () => {
    setDeckLoading(true)
    ApiGetRequest({
      route: "/tournament_player/tournament",
      query: { tournament_id: props.params.tournamentID },
      errorHandler: (err) => {
        setDeckLoading(false)
      },
      responseHandler: (res: { deck: Deck }) => {
        setDeckLoading(false)
        setDeck(res.deck)
      }
    })
  }

  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        <Header title={deck?.name || ""} endContent={<div className="text-gray-300 self-end font-bold">Deck</div>} />
        <div className="flex flex-col gap-2 mb-2">
          <MiniHeader title="Main Deck" />
          <MiniHeader title="Sideboard" />
          <MiniHeader title="Considering" />
        </div>
      </div>
    </Layout>
  )
}

type AddCardsModalProps = {
  isOpen: boolean
  closeFn: () => void
  refreshDecksFn: () => void
}

function AddCardsModal(props: AddCardsModalProps) {
  let [deckName, setDeckName] = useState<string>("")

  let createDeck = () => {
    // request
    props.refreshDecksFn()
  }

  return (
    <Modal
      isOpen={props.isOpen}
      placement="top-center"
    >
      <ModalContent>
        <ModalHeader className="flex flex-col text-white gap-1">Create new deck</ModalHeader>
        <ModalBody>
          <Input
            autoFocus
            className="text-white"
            label="Deck name"
            placeholder="Enter your deck's name"
            variant="bordered"
            onChange={(e) => setDeckName(e.target.value)}
          />
        </ModalBody>
        <ModalFooter>
          <Button color="danger" variant="flat" onPress={props.closeFn}>
            Cancel
          </Button>
          <Button color="success" onPress={() => { createDeck(); props.closeFn() }}>
            Create
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}