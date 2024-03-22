export type HeaderProps = {
  title: string
  endContent?: React.ReactElement | null
}

export function Header(props: HeaderProps) {
  return (
    <div className="my-8 flex flex-col">
      <div className="flex flex-row justify-between">
        <span className="self-center text-2xl font-semibold whitespace-nowrap text-white">{props.title}</span>
        {props.endContent}
      </div>
      <div className="w-full my-2 h-1 bg-white opacity-40" />
    </div>
  )
}

export function MiniHeader(props: HeaderProps) {
  return (
    <div className="flex flex-col w-full">
      <div className="flex flex-row justify-between">
        <span className="self-center text-md font-semibold whitespace-nowrap text-white">{props.title}</span>
        {props.endContent || null}
      </div>
      <div className="w-full my-2 h-1 bg-white opacity-20" />
    </div>
  )
}