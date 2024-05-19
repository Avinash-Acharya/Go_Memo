"use client";

import Link from "next/link";

export default function CreateButton() {
  return (
    <Link href="/memo-page/new">
      <button className=" absolute right-6 top-8 rounded-md bg-sky-800 p-1 text-sm  ">
        Create New
      </button>
    </Link>
  );
}

// const createData = async () => {
//   useEffect(() => {
//     createNewData();
//   }, []);
// }
// import { useEffect } from "react";
// import { createNewData } from "./action";
