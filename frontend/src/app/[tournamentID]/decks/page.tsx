"use client"

import { Header } from "@/components/header"
import Layout from "@/components/layout"
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests"
import { Deck } from "@/types/deck"
import { Button, Checkbox, Input, Link, Listbox, ListboxItem, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, Spinner, Textarea, useDisclosure } from "@nextui-org/react"
import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"
import { BsFillTrashFill } from "react-icons/bs";



export default function DecksPage(props: any) {
  let router = useRouter()
  let [isOpen, setIsOpen] = useState<boolean>(false)
  let [decks, setDecks] = useState<Deck[]>([])
  let [isLoading, setIsLoading] = useState<boolean>(false)
  let [error, setError] = useState<string>("")

  let refreshData = () => {
    setIsLoading(true)
    ApiGetRequest({
      route: "/deck/tournament_player",
      query: { tournament_id: props.params.tournamentID },
      errorHandler: (err) => {
        setError(err)
        setIsLoading(false)
      },
      responseHandler: (res: { decks: Deck[] }) => {
        setIsLoading(false)
        setDecks(res.decks)
      }
    })
  }

  useEffect(() => refreshData(), [props.params.tournamentID])

  return (
    <>
      <CreateDeckModal tournamentID={props.params.tournamentID} isOpen={isOpen} closeFn={() => setIsOpen(false)} refreshDecksFn={refreshData} />
      <Layout tournamentID={props.params.tournamentID}>
        <div className="mx-16 my-16">
          <Header title="Decks" endContent={<Button isIconOnly onClick={() => setIsOpen(true)} color="success">+</Button>} />
          {isLoading ? <div className="flex justify-center"> <Spinner /></div> :
            <div className="flex flex-col gap-2 mb-2">
              <div className="bg-gray-800 w-full border-small px-1 py-2 rounded-small border-default-200">
                <Listbox
                  emptyContent="No decks to show"
                >
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
          }
        </div>
      </Layout>
    </>
  )
}

type CreateDeckModalProps = {
  tournamentID: string
  isOpen: boolean
  closeFn: () => void
  refreshDecksFn: () => void
}

function CreateDeckModal(props: CreateDeckModalProps) {
  let [deckName, setDeckName] = useState<string>("")
  let [deckDescription, setDeckDescription] = useState<string>("")
  let [error, setError] = useState<string>("")
  let [isLoading, setIsLoading] = useState<boolean>(false)

  let sendCreateDeckRequest = () => {
    setError("")
    setIsLoading(true)
    ApiPostRequest({
      route: "/deck",
      body: {
        deck: { name: deckName, description: deckDescription },
        tournament_id: props.tournamentID
      },
      errorHandler: (err) => {
        setIsLoading(false)
        setError(err)
        props.refreshDecksFn()
      },
      responseHandler: () => {
        setIsLoading(false)
        props.closeFn()
        props.refreshDecksFn()
      }
    })
  }

  return (
    <Modal
      hideCloseButton
      onClose={props.closeFn}
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
            onValueChange={(value) => setDeckName(value)}
            isDisabled={isLoading}
          />
          <Textarea
            className="text-white"
            label="Deck description"
            placeholder=""
            variant="bordered"
            onValueChange={(value) => setDeckDescription(value)}
            isDisabled={isLoading}
          />
          <p className="text-sm font-light text-red-400 h-2">{error}</p>
        </ModalBody>
        <ModalFooter>
          <Button isDisabled={isLoading} color="danger" variant="flat" onPress={props.closeFn}>
            Cancel
          </Button>
          <Button isLoading={isLoading} color="success" onPress={sendCreateDeckRequest}>
            Create
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}