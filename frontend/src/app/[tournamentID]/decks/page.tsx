"use client"

import { Header } from "@/components/header"
import Layout from "@/components/layout"
import { Button, Checkbox, Input, Link, Listbox, ListboxItem, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, useDisclosure } from "@nextui-org/react"
import { useRouter } from "next/navigation"
import { useState } from "react"
import { BsFillTrashFill } from "react-icons/bs";



export default function DecksPage(props: any) {
  let router = useRouter()
  let [isOpen, setIsOpen] = useState<boolean>(false)

  let decks = [
    {
      name: "deck1",
      id: "1234"
    }
  ]

  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        <Header title="Decks" endContent={<Button onClick={() => setIsOpen(true)} color="success">+</Button>} />
        <div className="flex flex-col gap-2 mb-2">

          <CreateDeckModal isOpen={isOpen} closeFn={() => setIsOpen(false)} refreshDecksFn={() => console.log("refresh")} />
          <div className="bg-gray-800 w-full border-small px-1 py-2 rounded-small border-default-200">
            <Listbox>
              {
                decks.map((deck) =>
                  <ListboxItem
                    className="text-white"
                    onPress={() => router.push(`/${props.params.tournamentID}/decks/${deck.id}`)}
                    key={deck.name}
                  >
                    {deck.name}
                  </ListboxItem>
                )
              }
            </Listbox>
          </div>
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