export type HeaderProps = {
  title: string
}

export function Header(props: HeaderProps) {
  return (
    <div className="my-8">
      <span className="self-center text-2xl font-semibold whitespace-nowrap text-white">{props.title}</span>
      <div className="w-full my-2 h-1 bg-white opacity-20" />
    </div>
  )
}