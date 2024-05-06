import React from "react";
import ReactLoading from "react-loading";
import {
  Card,
  CardContent,
  CardDescription,
  CardTitle,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Upload } from "lucide-react";
import * as AlertDialog from "@/components/ui/alert-dialog";
import * as Table from "@/components/ui/table";
import { OpenFile, MekanoBilling } from "@/lib/wailsjs/go/main/App";
import { mekano } from "@/lib/wailsjs/go/models";

function Billing() {
  const [billFilePath, setBillFilePath] = React.useState("");
  const [extrasFilePath, setExtrasFilePath] = React.useState("");
  const [processing, setProcessing] = React.useState(false);
  const [result, setResult] = React.useState<mekano.BillingStatistics | null>(
    null
  );

  const OpenBillDialog = () => {
    OpenFile().then((path) => {
      setBillFilePath(path);
    });
  };

  const OpenExtraDialog = () => {
    OpenFile().then((path) => {
      setExtrasFilePath(path);
    });
  };

  const ProcessBilling = () => {
    setProcessing(true);
    MekanoBilling(billFilePath, extrasFilePath).then((result) => {
      setResult(result);
      setProcessing(false);
    });
  };

  const handleClear = () => {
    setBillFilePath("");
    setExtrasFilePath("");
    setResult(null);
  }
  return (
    <React.Fragment>
    <Card>
      <CardHeader>
        <CardTitle>Facturacion</CardTitle>
        <CardDescription>Carga el archivo de facturacion aqui.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-2">
        <div className="space-y-1">
          <Label className="font-semibold" htmlFor="pay_file">
            Archivo de Facturacion
          </Label>
          <div className="relative">
            <Upload
              onClick={OpenBillDialog}
              size={35}
              className="absolute flex items-center justify-center pl-3 bg-transparent text-black"
            />
            <Input placeholder={billFilePath} className="pl-28"></Input>
          </div>
        </div>
        <div className="space-y-1">
          <Label className="font-semibold" htmlFor="pay_file">
            Clientes con 2 o mas cuentas
          </Label>
          <div className="relative">
            <Upload
              size={35}
              className="absolute flex items-center justify-center pl-3 bg-transparent text-black"
              onClick={OpenExtraDialog}
            />

            <Input placeholder={extrasFilePath} className="pl-28"></Input>
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <Button
          className="w-full"
          onClick={ProcessBilling}
          disabled={processing}
        >
          {processing ? (
            <ReactLoading type={"spin"} height={20} width={20} />
          ) : (
            "Procesar"
          )}
        </Button>
      </CardFooter>
    </Card>
    {result !== null ? (
                    <AlertDialog.AlertDialog defaultOpen={true}>
                        <AlertDialog.AlertDialogContent>
                            <AlertDialog.AlertDialogHeader>
                                <AlertDialog.AlertDialogTitle>Resultados</AlertDialog.AlertDialogTitle>
                                <AlertDialog.AlertDialogDescription>
                                    <Table.Table>
                                        <Table.TableBody>
                                            <Table.TableRow>
                                                <Table.TableCell className="font-medium">Debito</Table.TableCell>
                                                <Table.TableCell className="text-right">${new Intl.NumberFormat().format(result.debito)}</Table.TableCell>
                                            </Table.TableRow>
                                            <Table.TableRow>
                                                <Table.TableCell className="font-medium">Credito</Table.TableCell>
                                                <Table.TableCell className="text-right">${new Intl.NumberFormat().format(result.credito)}</Table.TableCell>
                                            </Table.TableRow>
                                            <Table.TableRow>
                                                <Table.TableCell className="font-medium">Base</Table.TableCell>
                                                <Table.TableCell className="text-right">${new Intl.NumberFormat().format(result.base)}</Table.TableCell>
                                            </Table.TableRow>
                                        </Table.TableBody>
                                    </Table.Table>

                                </AlertDialog.AlertDialogDescription>
                            </AlertDialog.AlertDialogHeader>
                            <AlertDialog.AlertDialogFooter>
                                <AlertDialog.AlertDialogAction onClick={handleClear}>Continue</AlertDialog.AlertDialogAction>
                            </AlertDialog.AlertDialogFooter>
                        </AlertDialog.AlertDialogContent>
                    </AlertDialog.AlertDialog>

                ) : null}
    </React.Fragment>
  );
}

export default Billing;
