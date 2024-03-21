"use client"

import { Header, MiniHeader } from "@/components/header"
import Layout from "@/components/layout"
import { CreateTournamentPostModal } from "@/modals/create_tournament_post"
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests"
import { TournamentPlayer } from "@/types/tournamentPlayer"
import { TournamentPost } from "@/types/tournament_post"
import { Accordion, AccordionItem, Button, Spinner } from "@nextui-org/react"
import Image from "next/image"
import { useEffect, useState } from "react"

export default function TournamentHome(props: any) {
  let [createPostModalOpen, setCreatePostModalOpen] = useState<boolean>(false)
  let [tournamentPlayer, setTournamentPlayer] = useState<TournamentPlayer>()

  let refreshData = () => {
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

  useEffect(() => {
    refreshData()
  }, [])

  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        <Header title="Home" />
        <div className="flex flex-row w-full gap-8">
          <div className="flex flex-col max-w-[70%] w-full">
            <MiniHeader title="Posts" endContent={<Button size="sm" isIconOnly onClick={() => setCreatePostModalOpen(true)} color="success">+</Button>} />
            <CreateTournamentPostModal closeFn={() => setCreatePostModalOpen(false)} isOpen={createPostModalOpen} refreshFn={() => { }} tournamentID={props.params.tournamentID} />
            <TournamentPostsSection tournamentPlayer={tournamentPlayer} tournamentID={props.params.tournamentID} />
          </div>
          <MiniHeader title="Packs" />
        </div>
      </div>
    </Layout>
  )
}

type TournamentPostsSection = {
  tournamentID: string
  tournamentPlayer?: TournamentPlayer
}

function TournamentPostsSection(props: TournamentPostsSection) {
  let [posts, setPosts] = useState<TournamentPost[]>([])
  let [isLoading, setIsLoading] = useState<boolean>(true)

  let refreshData = () => {
    setIsLoading(true)
    ApiGetRequest({
      route: "/tournament_post",
      query: {
        tournament_id: props.tournamentID
      },
      errorHandler: () => { setIsLoading(false) },
      responseHandler: (res: { tournament_posts: TournamentPost[] }) => {
        setPosts(res.tournament_posts)
        setIsLoading(false)
      }
    })
  }

  let deletePost = (tournamentPostID: string) => {
    setIsLoading(true)
    ApiPostRequest({
      route: "/tournament_post/remove",
      query: {
        tournament_id: props.tournamentID
      },
      body: {
        tournament_post_id: tournamentPostID,
      },
      errorHandler: () => { setIsLoading(false) },
      responseHandler: () => {
        refreshData()
        setIsLoading(false)
      }
    })
  }

  useEffect(() => {
    refreshData()
  }, [props.tournamentID])

  if (!isLoading && !posts) {
    return (
      <div className="text-white bg-gray-800 p-4 rounded-lg">
        <div className="text-white font-thin italic mb-2">No posts to show</div>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-2">
      {isLoading ? <div className="flex justify-center"> <Spinner /></div> : posts.slice(0, 8).map((post) => (
        <div className="flex flex-col gap-2 bg-gray-800 p-4 rounded-xl">
          {post && props.tournamentPlayer && (props.tournamentPlayer.access_level == "al_administrator" || props.tournamentPlayer.access_level == "al_moderator") &&
            <Button isDisabled={isLoading} color="danger" variant="flat" onPress={() => deletePost(post.id)}>
              Delete post
            </Button>
          }
          <div className="text-lg font-bold mb-2 text-white">{post?.title || ""}</div>
          <Accordion className="text-white" isCompact>
            {post.blocks && post.blocks.map((block, index) => {
              return (
                <AccordionItem className="font-bold" key={index} title={block.title}>
                  <div className="flex flex-col gap-4">
                    {block.content.map((content) => {
                      switch (content.type) {
                        case "tbct_text":
                          return (
                            <div className="text-md font-thin">
                              {content.content}
                            </div>
                          )
                        case "tbct_image":
                          return (
                            <div className="flex justify-center w-full">
                              <Image className="rounded-2xl" src={content.content.startsWith("https://") ? content.content : ""} alt="image" width={300} height={300} />
                            </div>
                          )
                        default:
                          return ""
                      }
                    })}
                  </div>
                </AccordionItem>
              )
            })}
          </Accordion>
        </div>
      )
      )}
    </div >
  )
}