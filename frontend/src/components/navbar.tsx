import Image from "next/image"

export type NavigationTopbarProps = {
  toggleSidebarFn: any
}

export default function NavigationTopbar(props: NavigationTopbarProps) {
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
            <a href="#" className="flex mr-4">
              <Image src="/logo.png" className="mr-3 h-8" alt="WDML Logo" width={32} height={32} />
              <span className="self-center text-2xl font-semibold whitespace-nowrap text-white">WDML</span>
            </a>
          </div>
        </div>
      </nav >
    </aside >
  )
}