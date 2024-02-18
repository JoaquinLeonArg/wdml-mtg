"use client"

import NavigationTopbar from '@/components/navbar';
import NavigationSidebar from '@/components/sidebar';
import { useState } from 'react';

export default function Home() {
    let [sidebarVisible, setSidebarVisible] = useState<boolean>(true)
    return (
        <>
            <NavigationTopbar toggleSidebarFn={() => setSidebarVisible(!sidebarVisible)} />
            <NavigationSidebar visible={sidebarVisible} />
        </>
    )
}
