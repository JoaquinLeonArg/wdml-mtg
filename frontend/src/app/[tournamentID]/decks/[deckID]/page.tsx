"use client"

import { Header, MiniHeader } from "@/components/header"
import Layout from "@/components/layout"
import { Button, Checkbox, Input, Link, Listbox, ListboxItem, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, useDisclosure } from "@nextui-org/react"
import { useRouter } from "next/navigation"
import { useState } from "react"



export default function EditDeckPage(props: any) {
  let router = useRouter()
  let [isOpen, setIsOpen] = useState<boolean>(false)

  let deck = {
    name: "deck1",
    id: "1234",
    cards: [
      {
        image_url: "https://cards.scryfall.io/normal/front/0/1/0141312f-4b68-4c56-b1dc-5b7e6afbb96c.jpg?1627428247"
      }
    ]
  }


  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        <Header title={deck.name} endContent={<div className="text-gray-300 self-end font-bold">Deck</div>} />
        <div className="flex flex-col gap-2 mb-2">
          <MiniHeader title="Main Deck" />
        </div>
      </div>
    </Layout>
  )
}

type CreateDeckModalProps = {
  isOpen: boolean
  closeFn: () => void
  refreshDecksFn: () => void
}

function CreateDeckModal(props: CreateDeckModalProps) {
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