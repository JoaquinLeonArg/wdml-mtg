"use client"

import { Header } from "@/components/header"
import Layout from "@/components/layout"

export default function TournamentHome(props: any) {
  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        <Header title="Homepage under construction" />
      </div>
    </Layout>
  )
}