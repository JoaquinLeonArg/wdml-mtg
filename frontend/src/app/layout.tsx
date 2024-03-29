"use client"

// import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import "./mana.css"
import { NextUIProvider } from "@nextui-org/react";

const inter = Inter({ subsets: ["latin"] });

// export const metadata: Metadata = {
//   title: "WDML",
//   description: "A walk down memory lane - Magic the Gathering",
// };



export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`dark ${inter.className}`}>
        <NextUIProvider>
          {children}
        </NextUIProvider>
      </body>
    </html>
  );
}
