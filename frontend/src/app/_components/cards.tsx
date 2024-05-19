import Link from "next/link";

type propType = {
  key: number;
  id: number;
  head: string;
  body: string;
  // createdAt: Date;
};

export default function Cards(prop: propType) {
  return (
    <div>
      <Link href={`/memo-page/` + prop.id}>
        <div className="mx-2 mb-3 h-[100px] cursor-pointer rounded-lg bg-gradient-to-b from-slate-900 to-slate-800 p-3  shadow-none transition duration-300 hover:shadow hover:shadow-slate-500 sm:mx-auto sm:w-[600px] lg:w-[900px]  ">
          <h2 className="h-8 truncate text-lg text-slate-300">{prop.head}</h2>
          <p className=" truncate text-slate-400">{prop.body}</p>
          {/* <div className=" text-right text-sm text-gray-500">
            {prop.createdAt.toLocaleString()}
          </div> */}
        </div>
      </Link>
    </div>
  );
}
