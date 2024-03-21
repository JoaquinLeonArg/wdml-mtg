import { TournamentPostBlock } from "@/types/tournament_post"
import { Modal, ModalContent, ModalHeader, ModalBody, Input, Textarea, ModalFooter, Button, Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Checkbox } from "@nextui-org/react"
import { useState } from "react"
import Image from "next/image"
import { ErrorBoundary } from "react-error-boundary";
import { ApiPostRequest } from "@/requests/requests";

export type CreateTournamentPostModalProps = {
  tournamentID: string
  isOpen: boolean
  closeFn: () => void
  refreshFn: () => void
}

const contentTypes = {
  "tbct_text": "Text",
  "tbct_image": "Image"
}

export function CreateTournamentPostModal(props: CreateTournamentPostModalProps) {
  let [isLoading, setIsLoading] = useState<boolean>(false)
  let [postTitle, setPostTitle] = useState<string>("")
  let [blocks, setBlocks] = useState<TournamentPostBlock[]>([{ collapsable: false, content: [], title: "" }])

  let sendCreateRequest = () => {
    setIsLoading(true)
    ApiPostRequest({
      route: "/tournament_post",
      query: {
        tournament_id: props.tournamentID,
      },
      body: {
        tournament_post: {
          title: postTitle,
          blocks,
        }
      },
      errorHandler: () => {
        setIsLoading(false)
      },
      responseHandler: () => {
        setIsLoading(false)
        props.refreshFn()
        props.closeFn()
      }
    })
  }

  return (
    <Modal
      hideCloseButton
      onClose={() => { props.closeFn(); setBlocks([{ collapsable: false, content: [], title: "" }]) }}
      isOpen={props.isOpen}
      placement="top-center"
      size="4xl"
      scrollBehavior="outside"
    >
      <ModalContent>
        <ModalHeader className="flex flex-col text-white gap-1">Create new post</ModalHeader>
        <ModalBody>
          <Input
            autoFocus
            className="text-white"
            label="Post title"
            placeholder="Enter your post's title"
            variant="bordered"
            onValueChange={(value) => setPostTitle(value)}
            isDisabled={isLoading}
          />
          {
            blocks.map((block, blockIndex) => {
              return (
                <div key={block.title + blockIndex} className="flex flex-row w-full gap-2">
                  <div className="flex flex-col w-full gap-2">
                    <Input
                      className="text-white"
                      label="Section title"
                      placeholder="Enter the section's title"
                      variant="bordered"
                      onValueChange={
                        (value) => {
                          let newBlocks = [...blocks]
                          newBlocks[blockIndex].title = value
                          setBlocks(newBlocks)
                        }
                      }
                      isDisabled={isLoading}
                    />
                    <Checkbox>Collapsable</Checkbox>
                    {
                      block.content.map((content, contentIndex) => {
                        return (
                          <div key={content.content + contentIndex} className="flex flex-row w-full gap-2">
                            <Dropdown>
                              <DropdownTrigger>
                                <Button
                                  className="h-full"
                                  variant="bordered"
                                >
                                  {contentTypes[content.type]}
                                </Button>
                              </DropdownTrigger>
                              <DropdownMenu aria-label="Static Actions">
                                <DropdownItem onClick={() => {
                                  let newBlocks = [...blocks]
                                  newBlocks[blockIndex].content[contentIndex].type = "tbct_text"
                                  setBlocks(newBlocks)
                                }} key="tbct_text" className="text-white">Text</DropdownItem>
                                <DropdownItem onClick={() => {
                                  let newBlocks = [...blocks]
                                  newBlocks[blockIndex].content[contentIndex].type = "tbct_image"
                                  setBlocks(newBlocks)
                                }} key="tbct_image" className="text-white">Image</DropdownItem>
                              </DropdownMenu>
                            </Dropdown>
                            {
                              content.type == "tbct_text" &&
                              <Textarea
                                className="text-white"
                                label="Text"
                                placeholder=""
                                variant="bordered"
                                onValueChange={
                                  (value) => {
                                    let newBlocks = [...blocks]
                                    newBlocks[blockIndex].content[contentIndex].content = value
                                    setBlocks(newBlocks)
                                  }
                                }
                                isDisabled={isLoading}
                              />
                            }
                            {
                              content.type == "tbct_image" &&
                              <div className="flex flex-row w-full gap-2 items-center">
                                <Input
                                  className="text-white"
                                  label="URL"
                                  placeholder="Enter the image url"
                                  variant="bordered"
                                  onValueChange={
                                    (value) => {
                                      let newBlocks = [...blocks]
                                      newBlocks[blockIndex].content[contentIndex].content = value
                                      setBlocks(newBlocks)
                                    }
                                  }
                                  isDisabled={isLoading}
                                />
                                <Image
                                  className="rounded-md"
                                  src={blocks[blockIndex].content[contentIndex].content.startsWith("https://cards.scryfall.io") ? blocks[blockIndex].content[contentIndex].content : ""}
                                  alt="preview" width={32} height={32} />
                              </div>
                            }
                            <Button onClick={() => {
                              let newBlocks = [...blocks]
                              newBlocks[blockIndex].content.splice(contentIndex, 1)
                              setBlocks(newBlocks)
                            }} color="danger" isIconOnly className="h-[100]">X</Button>
                          </div>
                        )
                      })
                    }
                    <Button onClick={() => {
                      let newBlocks = [...blocks]
                      newBlocks[blockIndex].content.push({
                        content: "",
                        extra_data: {},
                        type: "tbct_text"
                      })
                      setBlocks(newBlocks)
                    }} color="success" className="w-full">Add content</Button>
                  </div>
                  <Button onClick={() => {
                    let newBlocks = [...blocks]
                    newBlocks.splice(blockIndex, 1)
                    setBlocks(newBlocks)
                  }} color="danger" isIconOnly className="h-[100]">X</Button>
                </div>
              )
            })
          }
          <Button onClick={() => {
            let newBlocks = [...blocks]
            newBlocks.push({
              collapsable: false,
              content: [],
              title: ""
            })
            setBlocks(newBlocks)
          }} color="success" className="w-full">Add block</Button>
          {/* <p className="text-sm font-light text-red-400 h-2">{error}</p> */}
        </ModalBody>
        <ModalFooter>
          <Button isDisabled={isLoading} color="danger" variant="flat" onPress={props.closeFn}>
            Cancel
          </Button>
          <Button isLoading={isLoading} color="success" onPress={sendCreateRequest}>
            Create
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal >
  )
}
