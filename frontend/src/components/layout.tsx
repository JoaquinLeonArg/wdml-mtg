"use client"

import NavigationTopbar from '@/components/navbar';
import NavigationSidebar from '@/components/sidebar';
import { useState } from 'react';

export type LayoutProps = React.PropsWithChildren & {
    tournamentID: string
}

export default function Layout(props: LayoutProps) {
    let [sidebarVisible, setSidebarVisible] = useState<boolean>(true)
    return (
        <>
            <NavigationTopbar toggleSidebarFn={() => setSidebarVisible(!sidebarVisible)} />
            <NavigationSidebar tournamentID={props.tournamentID} visible={sidebarVisible} />
            <main className={`${sidebarVisible && "ml-64"} p-8`}>
                {props.children}
            </main>
        </>
    )
}
