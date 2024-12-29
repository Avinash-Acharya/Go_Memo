"use client";

import axios from "axios";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000/api/go/users";

type dataType = {
  id: number;
  head: string;
  body: string;
  createdAt: string;
};

export default function Memo({ params }: { params: { id: string } }) {
  const [head, setHead] = useState<string>("");
  const [body, setBody] = useState<string>("");
  const [isNew, setIsNew] = useState<boolean>(true);
  const router = useRouter();

  const idNo = params.id;

  useEffect(() => {
    const fetchData = async () => {
      if (idNo !== "new") {
        try {
          const response = await axios.get(URL + "/" + idNo);
          const dataById = response.data as dataType;
          if (dataById) {
            setHead(dataById.head);
            setBody(dataById.body);
          }
          setIsNew(false);
        } catch (err) {
          console.error(err);
        }
      }
    };

    void fetchData();
  }, [idNo]);

  const handleDelete = async () => {
    try {
      await axios.delete(URL + "/" + idNo);
    } catch (err) {
      console.log(err);
    }
    router.push("/"); // Navigate to the home page
    router.refresh();
  };

  const handleSave = async () => {
    try {
      if (isNew) {
        await axios.post(URL, {
          head: head,
          body: body,
          createdAt: new Date().toLocaleString(),
        });
      } else {
        await axios.put(URL + "/" + idNo, {
          head: head,
          body: body,
          createdAt: new Date().toLocaleString(),
        });
      }
      router.push("/"); // Navigate to the home page
      router.refresh();
    } catch (err) {
      console.log(err);
    }
  };

  return (
    <div className="flex min-h-screen flex-col bg-gradient-to-b from-[#1a0b2e] to-[#000000] text-white">
      <header className=" basis-1/12 py-6 text-center">
        <h1 className="text-2xl">
          <Link href={"/"}>
            <span className="italic">M</span>E
            <span className="underline">M</span>@
          </Link>
        </h1>
      </header>
      <main className="mx-2 mb-6 flex flex-grow basis-11/12 flex-col overflow-hidden rounded-lg border-4 border-teal-800 border-opacity-10 sm:mx-32">
        <input
          type="text"
          placeholder="Heading..."
          value={head}
          onChange={(e) => setHead(e.target.value)}
          className="h-14 w-full truncate bg-slate-800 bg-opacity-30 indent-5 text-xl outline-none"
        />
        <textarea
          placeholder="Type here...."
          value={body}
          onChange={(e) => setBody(e.target.value)}
          className="text-md h-full w-full flex-grow resize-none scroll-p-4 bg-slate-900 bg-opacity-25 p-6 outline-none "
        />
        <div className="flex">
          <button
            onClick={handleDelete}
            className="flex-grow bg-slate-900 bg-opacity-50 p-2 text-stone-300 transition duration-200 hover:bg-red-900"
          >
            Delete
          </button>
          <button
            onClick={handleSave}
            className="flex-grow bg-slate-900 bg-opacity-50 p-2 text-stone-300 transition duration-200 hover:bg-lime-900"
          >
            Save
          </button>
        </div>
      </main>
    </div>
  );
}
