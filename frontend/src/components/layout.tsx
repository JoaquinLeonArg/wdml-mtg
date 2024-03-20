"use client"

import NavigationTopbar from '@/components/navbar';
import NavigationSidebar from '@/components/sidebar';
import { ApiGetRequest } from '@/requests/requests';
import { Tournament } from '@/types/tournament';
import { useEffect, useState } from 'react';

export type LayoutProps = React.PropsWithChildren & {
  tournamentID: string
}

export default function Layout(props: LayoutProps) {
  let [sidebarVisible, setSidebarVisible] = useState<boolean>(true)
  let [tournaments, setTournaments] = useState<Tournament[]>([])

  let refreshData = () => {
    ApiGetRequest({
      route: "/tournament/user",
      query: {
        tournament_id: props.tournamentID,
      },
      errorHandler: (err) => { },
      responseHandler: (res: { tournaments: Tournament[] }) => {
        setTournaments(res.tournaments)
      }
    })
  }

  useEffect(() => {
    refreshData()
  }, [props.tournamentID])

  return (
    <>
      <NavigationTopbar tournamentID={props.tournamentID} tournaments={tournaments} toggleSidebarFn={() => setSidebarVisible(!sidebarVisible)} />
      <NavigationSidebar tournamentID={props.tournamentID} visible={sidebarVisible} />
      <main className={`${sidebarVisible && "ml-64"} p-8 min-h-[100vh]`}>
        {props.children}
      </main>
    </>
  )
}
