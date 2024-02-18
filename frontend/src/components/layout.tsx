"use client"

import NavigationTopbar from '@/components/navbar';
import NavigationSidebar from '@/components/sidebar';
import { useState } from 'react';

export default function Layout(props: React.PropsWithChildren) {
    let [sidebarVisible, setSidebarVisible] = useState<boolean>(true)
    return (
        <>
            <NavigationTopbar toggleSidebarFn={() => setSidebarVisible(!sidebarVisible)} />
            <NavigationSidebar visible={sidebarVisible} />
            <main className={`${sidebarVisible && "ml-64"} p-8`}>
                {props.children}
            </main>
        </>
    )
}
