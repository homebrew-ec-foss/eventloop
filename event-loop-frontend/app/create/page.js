"use client";

import "@/app/globals.css";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/ui/use-toast";
import * as React from "react";

import { useState } from "react";
import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import Navbar from "@/components/navbar";
import Access from "@/components/access";


export default function FileUpload() {
    // const [file, setFile] = useState(null);
    // const [status, setStatus] = useState("");
    // const [content, setContent] = useState("");
    const [jsonResponse, setJsonResponse] = useState(null);
    const { toast } = useToast();

    const newHandleSubmit = async (event) => {
        event.preventDefault();
        const eventName = document.getElementById("event-name").value;
        // const eventDate = document.getElementById("event-date").value;
        const file = document.querySelector('input[type="file"]').files[0];

        const formData = new FormData();
        formData.append("file", file);

        try {
            const userData = localStorage.getItem("google-oauth");
            if (!userData) {
                console.error("Local storage is empty");
            }

            const parsedUserData = JSON.parse(userData);

            formData.append("sub", parsedUserData.sub);

            const response = await fetch(`${process.env.GO_BACKEND_URL}/${parsedUserData.userRole}/create`, {
                method: "POST",
                body: formData,
            });

            const resp = await response.json();
            console.log(resp);

            if (response.ok) {
                console.log("Data is der");
                setJsonResponse(resp.data);
                console.log(resp.data);
            } else {
                console.log("Failed to upload file");
            }
        } catch (error) {
            console.log("Error: " + error.message);
        }
    };

    return (
        <main className="flex min-h-screen flex-col p-5 md:p-28 gap-4">
            <Navbar />

            <Access userRole={["admin", "operator"]} />

            {/*
                TODO:   Make all fields as "required"
                        removed for testing
            */}
            <Card className="max-[300px]">
                <CardHeader>
                    <CardTitle>Create Event</CardTitle>
                    <CardDescription>
                        Please enter relevant inputs
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={newHandleSubmit}>
                        <div className="flex flex-col gap-4">
                            <div className="flex flex-col space-y-1.5">
                                <Label>Event Name</Label>
                                <Input
                                    id="event-name"
                                    type="text"
                                    placeholder="Enter your event name"
                                />
                            </div>
                            <div className="flex flex-col space-y-1.5">
                                <Label>Event Date</Label>
                                <Input
                                    id="event-date"
                                    type="date"
                                    placeholder="Enter event date"
                                />
                            </div>
                            <div className="flex flex-col space-y-1.5">
                                <Label>Registrants CSV</Label>
                                <Input
                                    type="file"
                                    accept=".csv"
                                    multiple={false}
                                />
                            </div>
                            <div className="flex flex-col space-y-1.5">
                                <Button
                                    type="submit"
                                    onClick={() => {
                                        toast({
                                            title: "Event Created",
                                            description:
                                                "Event has been created successfully",
                                            type: "success",
                                        });
                                    }}
                                >
                                    Submit
                                </Button>
                            </div>
                        </div>
                    </form>
                </CardContent>
            </Card>

            {/*
                TODO:   Provide admin with stats regarding the form
                        before displaying the entire structure
            */}

            {jsonResponse && (
                <Table className="">
                    {/* fill table with setJsonResponse */}
                    <TableCaption>Verified CSV form</TableCaption>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Name</TableHead>
                            <TableHead>Phone</TableHead>
                            <TableHead>Email</TableHead>
                            <TableHead>College</TableHead>
                            <TableHead>Branch</TableHead>
                            <TableHead>PES Hostel</TableHead>
                            <TableHead>Team Name</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {jsonResponse &&
                            jsonResponse.map((row, index) => (
                                <TableRow key={index}>
                                    <TableCell>{row.name}</TableCell>
                                    <TableCell>{row.phone}</TableCell>
                                    <TableCell>{row.email}</TableCell>
                                    <TableCell>{row.college}</TableCell>
                                    <TableCell>{row.branch}</TableCell>
                                    {row.pesHostel != "" ? (
                                        <TableCell>{row.pesHostel}</TableCell>
                                    ) : (
                                        <TableCell>nil</TableCell>
                                    )}
                                    <TableCell>{row.team}</TableCell>
                                </TableRow>
                            ))}
                    </TableBody>
                </Table>
            )}
        </main>
    );
}
