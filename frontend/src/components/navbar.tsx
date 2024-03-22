import { Tournament } from "@/types/tournament"
import { Dropdown, DropdownTrigger, Button, DropdownMenu, DropdownItem } from "@nextui-org/react"
import Image from "next/image"
import { useRouter } from "next/navigation"
import { BsShareFill } from "react-icons/bs"

export type NavigationTopbarProps = {
  toggleSidebarFn: any
  tournaments?: Tournament[]
  tournamentID: string
}

export default function NavigationTopbar(props: NavigationTopbarProps) {
  let router = useRouter()

  return (
    <aside className="antialiased fixed top-0 w-full z-[200]">
      <nav className="border-gray-200 px-4 lg:px-6 py-2.5 bg-gray-600">
        <div className="flex flex-wrap justify-between items-center">
          <div className="flex justify-start items-center">
            <button onClick={() => props.toggleSidebarFn()} id="toggleSidebar" aria-expanded="true" aria-controls="sidebar" className="hidden p-2 mr-3 rounded cursor-pointer lg:inline text-gray-400 hover:text-white hover:bg-gray-700">
              <svg className="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 12"> <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M1 1h14M1 6h14M1 11h7" /> </svg>
            </button>
            <button onClick={() => props.toggleSidebarFn()} aria-expanded="true" aria-controls="sidebar" className="p-2 mr-4 rounded-lg cursor-pointer lg:hidden  focus:bg-gray-700 focus:ring-2 focus:ring-gray-700 text-gray-400 hover:bg-gray-700 hover:text-white">
              <svg className="w-[18px] h-[18px]" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 17 14"><path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M1 1h15M1 7h15M1 13h15" /></svg>
              <span className="sr-only">Toggle sidebar</span>
            </button>
            <a href="#" className="flex gap-12">
              <div className="flex flex-row">
                <Image src="/logo.png" className="mr-3 h-8" alt="TA Logo" width={32} height={32} />
                <span className="self-center text-2xl font-semibold whitespace-nowrap text-white">Tolarian Archives</span>
              </div>
              <div className="flex flex-row items-center gap-1">
                {props.tournaments && props.tournaments.length > 0 &&
                  <>
                    <Dropdown>
                      <DropdownTrigger>
                        <Button
                          variant="bordered"
                          size="sm"
                          className="w-48"
                        >
                          {props.tournaments.find((tournament) => tournament.id == props.tournamentID)?.name || ""}
                        </Button>
                      </DropdownTrigger>
                      <DropdownMenu aria-label="Tournament Selection" items={props.tournaments.concat({ id: "!", name: "Join or create", invite_code: "", description: "" })}>
                        {(item) => {
                          if (item.id == "!") {
                            return (
                              <DropdownItem
                                key="join"
                                color="default"
                                className="text-white"
                                onClick={() => router.push("/join")}
                              >
                                {item.name}
                              </DropdownItem>

                            )
                          }
                          return (
                            <DropdownItem
                              key={item.id}
                              color="default"
                              showDivider={props.tournaments && item.id == props.tournaments[props.tournaments.length - 1].id}
                              className="text-white"
                              onClick={() => {
                                if (item.id != props.tournamentID) {
                                  router.push("/" + item.id)
                                }
                              }}
                            >
                              {item.name}
                            </DropdownItem>
                          )
                        }}
                      </DropdownMenu>
                    </Dropdown>
                    <Button isIconOnly variant="bordered" size="sm"
                      onClick={() => {
                        if (props.tournaments)
                          navigator.clipboard.writeText(props.tournaments.find(tournament => tournament.id == props.tournamentID)?.invite_code || "")
                      }}
                    >
                      <BsShareFill />
                    </Button>
                  </>
                }
              </div>
            </a>
          </div>
        </div>
      </nav >
    </aside >
  )
}