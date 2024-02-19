export type PackListProps = {
  packs: Pack[]
}

export type Pack = {
  setName: string
  setCode: string
  count: number
}

export function PackList(props: PackListProps) {
  return (
    <div className="">
      {
        props.packs.map((pack: Pack, i: number) => {
          return (
            <BoosterPack key={i} pack={pack} />
          )
        })
      }
    </div>
  )
}

type BoosterPackProps = {
  pack: Pack
}

function BoosterPack(props: BoosterPackProps) {
  return (
    <div className="w-72 h-16 bg-primary-400 rounded-lg my-2 pl-4 text-base flex flex-row items-center justify-between">
      <div className="flex flex-col">
        <span className="mt-2">
          {props.pack.setName}
        </span>
        <span className="font-bold text-2xl">
          {props.pack.setCode}
        </span>

      </div>
      <svg className="mr-4 w-8 h-8" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z"></path></svg>
    </div>
  )
}