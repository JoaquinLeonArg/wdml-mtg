import { useRouter } from "next/navigation"
import { useState } from "react"
import { BsMegaphoneFill, BsBoxSeamFill, BsCollectionFill, BsDiagram3Fill, BsEmojiSunglassesFill, BsFillGearFill } from "react-icons/bs";

export type NavigationSidebarProps = {
    visible: boolean
}

export default function NavigationSidebar(props: NavigationSidebarProps) {
    return (
        <div className={`${!props.visible && "hidden"}`}>
            <aside id="default-sidebar" className="fixed top-17 left-0 z-40 w-64 h-screen transition-transform sm:translate-x-0" aria-label="Sidenav">
                <div className="overflow-y-auto py-5 px-3 h-full bg-gray-800 border-gray-700">
                    <ul className="space-y-2">
                        <SidebarButton name="News" ><BsMegaphoneFill className="w-6 h-6" /></SidebarButton>
                        <SidebarButton name="Packs" ><BsBoxSeamFill className="w-6 h-6" /></SidebarButton>
                        <SidebarButton name="Cards" items={[{ name: "Collection", href: "#" }, { name: "Decks", href: "#" }]} ><BsCollectionFill className="w-6 h-6" /></SidebarButton>
                        <SidebarButton name="Matches" ><BsDiagram3Fill className="w-6 h-6" /></SidebarButton>
                        <SidebarButton name="Players" ><BsEmojiSunglassesFill className="w-6 h-6" /></SidebarButton>
                        <SidebarButton name="Settings" ><BsFillGearFill className="w-6 h-6" /></SidebarButton>
                    </ul>

                </div>

            </aside>
        </div>
    )
}

type SidebarButtonProps = React.PropsWithChildren & {
    name: string
    href?: string
    items?: SidebarButtonItem[]
}

type SidebarButtonItem = {
    name: string
    href: string
}

function SidebarButton(props: SidebarButtonProps) {
    let [showItems, setShowItems] = useState<boolean>(false)
    let router = useRouter()
    return (
        <li>
            <button type="button" onClick={() => { if (props.items?.length) setShowItems(!showItems); else router.push(props.href || "#") }} className="flex items-center p-2 w-full text-base font-normal rounded-lg transition duration-75 group text-white hover:bg-gray-700" aria-controls={props.items?.length ? `dropdown-${props.name}` : ""} data-collapse-toggle={props.items?.length ? `dropdown-${props.name}` : ""}>
                {props.children}
                <span className="flex-1 ml-3 text-left whitespace-nowrap">{props.name}</span>
                {/* Arrow */}
                {(props.items && props.items.length > 0) && <svg aria-hidden="true" className="w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>}
            </button>
            <ul id="dropdown-pages" className={`${!showItems && "hidden"} py-2 space-y-2`}>
                {
                    props.items?.map((item: SidebarButtonItem, i: number) => {
                        return (
                            <li key={i}>
                                <a href={item.href} className="flex items-center p-2 pl-6 w-full text-sm font-normal rounded-lg transition duration-75 group text-white hover:bg-gray-700">{item.name}</a>
                            </li>
                        )
                    })
                }
            </ul>
        </li>
    )
}